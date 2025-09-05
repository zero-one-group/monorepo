package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zero-one-group/go-modulith/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func NewDatabase(cfg *config.Config) (*DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	if cfg.IsDevelopment() {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnectionMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnectionMaxIdleTime)

	slog.Info("Database connected successfully", 
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"database", cfg.Database.DBName)

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *DB) Health(ctx context.Context) error {
	ctx, span := otel.Tracer("database").Start(ctx, "health_check")
	defer span.End()

	sqlDB, err := db.DB.DB()
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return err
	}

	span.SetAttributes(attribute.Bool("healthy", true))
	return nil
}

func (db *DB) WithContext(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx)
}

func (db *DB) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	ctx, span := otel.Tracer("database").Start(ctx, "transaction")
	defer span.End()

	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("transaction.rolled_back", true))
			return err
		}
		span.SetAttributes(attribute.Bool("transaction.committed", true))
		return nil
	})
}