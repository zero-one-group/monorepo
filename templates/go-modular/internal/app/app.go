package app

import (
	"log/slog"

	"{{ package_name | kebab_case }}/internal/config"
	modAuth "{{ package_name | kebab_case }}/modules/auth"
	modUser "{{ package_name | kebab_case }}/modules/user"

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
