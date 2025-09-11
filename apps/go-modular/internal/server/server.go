package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-modular/internal/adapter"
	"go-modular/internal/config"
	"go-modular/internal/notification"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	appMiddleware "go-modular/internal/middleware"
	templateFS "go-modular/templates"
)

// HTTPServer is the main HTTP server struct.
// Logger and Tracer are injected for observability.
type HTTPServer struct {
	httpAddr string
	logger   *slog.Logger
}

func NewHTTPServer(httpAddr string, logger *slog.Logger) *HTTPServer {
	return &HTTPServer{
		httpAddr: httpAddr,
		logger:   logger,
	}
}

func (s *HTTPServer) Start() error {
	cfg := config.Get()

	// Initialize Postgres database connection with retry mechanism
	pg, err := s.initializeDatabase(cfg)
	if err != nil {
		s.logger.Error("Failed to connect to Postgres database", "err", err)
		os.Exit(1)
	}
	// ensure DB pool closed on return (also closed during graceful shutdown below)
	defer pg.Close()

	var mailer *notification.Mailer
	s.logger.Info("Initializing SMTP mailer service")
	m, err := notification.NewMailer(notification.MailerOptions{
		SMTPHost:     cfg.Mailer.SMTPHost,
		SMTPPort:     cfg.Mailer.SMTPPort,
		SMTPUsername: cfg.Mailer.SMTPUsername,
		SMTPPassword: cfg.Mailer.SMTPPassword,
		FromName:     cfg.Mailer.SenderName,
		FromAddress:  cfg.Mailer.SenderEmail,
		TemplateFS:   templateFS.TemplateDir,
		Logger:       s.logger,
	})
	if err != nil {
		s.logger.Info("Mailer service not configured or failed to initialize, continuing without mailer", "err", err)
	} else {
		mailer = m
		s.logger.Info("Mailer service initialized", "host", cfg.Mailer.SMTPHost, "port", cfg.Mailer.SMTPPort)
	}

	e := echo.New() // Create Echo instance
	e.Logger.SetLevel(cfg.GetEchoLogLevel())
	e.HideBanner = true
	e.HidePort = true

	// Register global middlewares
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		// Use the RecoverConfig.LogErrorFunc signature: func(c echo.Context, err error, stack []byte) error
		// Integrate with slog logger and return the error so the centralized HTTPErrorHandler still runs.
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			req := c.Request()
			s.logger.Error("Recovered from panic", "method", req.Method, "path", req.URL.Path, "error", err)
			// return the error to allow downstream error handler to process it
			return err
		},
	}))

	e.Use(appMiddleware.RequestIDMiddleware())
	e.Use(appMiddleware.SecurityHeadersMiddleware())         // Globally enable security headers
	e.Use(appMiddleware.TimeoutMiddleware(time.Second * 30)) // Maximum request timeout: 30s
	e.Use(appMiddleware.LoggerMiddleware(s.logger))

	// Register modules and their route, put inside helper to make Start tidy
	if err := s.registerModules(cfg, pg, mailer, e); err != nil {
		s.logger.Error("failed to register modules", "err", err)
		// decide whether to exit or continue; here we exit as server without routes is useless
		return err
	}

	s.logger.Info("Starting HTTP server", "addr", s.httpAddr)

	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in background
	serverErrCh := make(chan error, 1)
	go func() {
		if err := e.Start(s.httpAddr); err != nil && err != http.ErrServerClosed {
			serverErrCh <- fmt.Errorf("server start failed: %w", err)
		}
		// If server ended without error, notify
		close(serverErrCh)
	}()

	// Wait for signal or server start error
	select {
	case <-ctx.Done():
		// received shutdown signal
		s.logger.Info("Shutdown signal received, shutting down HTTP server")
	case err := <-serverErrCh:
		if err != nil {
			// server failed to start or crashed
			s.logger.Error("HTTP server error", "err", err)
			// proceed to shutdown resources anyway
		}
	}

	// Shutdown timeout set to 10s
	shutdownTimeout := 10 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server gracefully
	if err := e.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("failed to shutdown HTTP server gracefully", "err", err)
	}

	// Close DB pool
	s.logger.Info("Closing database connections")
	pg.Close()

	// Terminate mailer if initialized. Try Shutdown(ctx) first, then Close()
	if mailer != nil {
		s.logger.Info("Shutting down mailer")
		if shutdowner, ok := any(mailer).(interface{ Shutdown(context.Context) error }); ok {
			if err := shutdowner.Shutdown(shutdownCtx); err != nil {
				s.logger.Error("mailer shutdown error", "err", err)
			}
		} else if closer, ok := any(mailer).(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				s.logger.Error("mailer close error", "err", err)
			}
		}
	}

	s.logger.Info("Shutdown complete")
	return nil
}

// Initialize PostgreSQL database connection with retry mechanism
func (s *HTTPServer) initializeDatabase(cfg *config.Config) (*adapter.PostgresDB, error) {
	const baseDelay = 2 * time.Second
	const maxDelay = 30 * time.Second
	const defaultMaxRetries = 5

	maxRetries := defaultMaxRetries
	if cfg.Database.PgMaxRetries == -1 {
		maxRetries = -1
	} else if cfg.Database.PgMaxRetries > 0 {
		maxRetries = cfg.Database.PgMaxRetries
	}

	s.logger.Info("Initializing database connection", "max_retries", maxRetries)

	var lastErr error
	attempt := 1

	for {
		pgCfg := adapter.PostgresConfig{URL: cfg.GetDatabaseURL()}
		pg, err := adapter.NewPostgres(pgCfg)
		if err == nil {
			// verify connection with Ping and timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			pingErr := pg.Ping(ctx)
			cancel()

			if pingErr == nil {
				s.logger.Info("Database connection established", "attempt", attempt)
				return pg, nil
			}

			// ping failed, close pool and treat as error
			pg.Close()
			err = fmt.Errorf("ping failed: %w", pingErr)
		}

		lastErr = err
		s.logger.Warn("Database connection attempt failed", "attempt", attempt, "err", lastErr)

		// If not infinite and we've reached max, stop retrying
		if maxRetries != -1 && attempt >= maxRetries {
			break
		}

		// Exponential-ish backoff bounded by maxDelay
		delay := min(baseDelay*time.Duration(attempt), maxDelay)

		s.logger.Info("Retrying database connection", "next_try_in", delay, "attempt", attempt+1)
		time.Sleep(delay)
		attempt++
	}

	return nil, fmt.Errorf("failed to establish database connection after %d attempts: %w", attempt, lastErr)
}

// helper min function
func min(x, y time.Duration) time.Duration {
	if x < y {
		return x
	}
	return y
}
