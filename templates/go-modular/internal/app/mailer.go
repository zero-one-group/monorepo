package app

import (
	"log/slog"

	"{{ package_name | kebab_case }}/internal/config"
	"{{ package_name | kebab_case }}/internal/notification"
	templateFS "{{ package_name | kebab_case }}/templates"

	"go.uber.org/fx"
)

var mailerModule = fx.Module("mailer",
	fx.Provide(newMailer),
)

// Soft-fail: misconfigured SMTP yields nil mailer; Auth still starts.
func newMailer(cfg *config.Config, log *slog.Logger) *notification.Mailer {
	log.Info("Initializing SMTP mailer service")
	m, err := notification.NewMailer(notification.MailerOptions{
		SMTPHost:     cfg.Mailer.SMTPHost,
		SMTPPort:     cfg.Mailer.SMTPPort,
		SMTPUsername: cfg.Mailer.SMTPUsername,
		SMTPPassword: cfg.Mailer.SMTPPassword,
		FromName:     cfg.Mailer.SenderName,
		FromAddress:  cfg.Mailer.SenderEmail,
		TemplateFS:   templateFS.TemplateDir,
		Logger:       log,
	})
	if err != nil {
		log.Info("Mailer service not configured or failed to initialize, continuing without mailer", "err", err)
		return nil
	}
	log.Info("Mailer service initialized", "host", cfg.Mailer.SMTPHost, "port", cfg.Mailer.SMTPPort)
	return m
}
