package app

import (
	"log/slog"

	"go-modular/internal/config"
	modAuth "go-modular/modules/auth"
	modUser "go-modular/modules/user"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// New builds the serve-path fx application.
func New(cfg *config.Config, log *slog.Logger) *fx.App {
	opts := []fx.Option{
		fx.Supply(cfg, log),
		postgresModule,
		mailerModule,
		httpModule,
		modUser.Module,
		modAuth.Module,
	}
	if cfg.IsProduction() {
		opts = append([]fx.Option{fx.NopLogger}, opts...)
	} else {
		opts = append([]fx.Option{
			fx.WithLogger(func(l *slog.Logger) fxevent.Logger {
				return &fxevent.SlogLogger{Logger: l}
			}),
		}, opts...)
	}
	return fx.New(opts...)
}
