package config

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// Validates critical configuration values
func validateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	var errs []string

	// App mode
	mode := strings.ToLower(strings.TrimSpace(config.App.Mode))
	if mode == "" {
		errs = append(errs, "app mode is required")
	} else {
		validModes := []string{"development", "production", "staging", "test"}
		if !slices.Contains(validModes, mode) {
			errs = append(errs, fmt.Sprintf("invalid app mode: %q (valid: %v)", config.App.Mode, validModes))
		}
	}

	// JWT algorithm and secret
	alg := strings.ToUpper(strings.TrimSpace(string(config.App.JWTAlgorithm)))
	if alg != "" {
		validAlgs := []string{string(JWTAlgorithmHS256), string(JWTAlgorithmRS256)}
		if !slices.Contains(validAlgs, alg) {
			errs = append(errs, fmt.Sprintf("invalid JWT algorithm: %q (valid: %v)", alg, validAlgs))
		}
	}
	secret := strings.TrimSpace(config.App.JWTSecretKey)
	if strings.EqualFold(mode, "production") {
		if secret == "" || secret == "_THIS_IS_DEFAULT_JWT_SECRET_KEY_" {
			errs = append(errs, "JWT secret key must be set in production")
		}
	}

	// Server port
	if config.App.ServerPort <= 0 || config.App.ServerPort > 65535 {
		errs = append(errs, fmt.Sprintf("invalid server port: %d (must be 1-65535)", config.App.ServerPort))
	}

	// Database URL
	dbURL := strings.TrimSpace(config.Database.PostgresURL)
	if dbURL == "" {
		errs = append(errs, "database URL is required")
	} else {
		if _, err := url.Parse(dbURL); err != nil {
			errs = append(errs, fmt.Sprintf("invalid database URL: %v", err))
		}
	}

	// Postgres pool size / retries
	if config.Database.PgMaxPoolSize <= 0 {
		errs = append(errs, fmt.Sprintf("pg max pool size must be > 0 (got %d)", config.Database.PgMaxPoolSize))
	}
	if config.Database.PgMaxRetries < 0 {
		errs = append(errs, fmt.Sprintf("pg max retries must be >= 0 (got %d)", config.Database.PgMaxRetries))
	}

	// Mailer
	mh := strings.TrimSpace(config.Mailer.SMTPHost)
	if mh != "" {
		if config.Mailer.SMTPPort <= 0 || config.Mailer.SMTPPort > 65535 {
			errs = append(errs, fmt.Sprintf("invalid mailer SMTP port: %d (must be 1-65535)", config.Mailer.SMTPPort))
		}
	} else {
		// If host is blank but port set, warn/error
		if config.Mailer.SMTPPort > 0 {
			errs = append(errs, "mailer SMTP host is empty but SMTP port is set")
		}
	}

	// File store / S3
	s3ep := strings.TrimSpace(config.FileStore.S3Endpoint)
	if s3ep != "" {
		if strings.TrimSpace(config.FileStore.S3BucketName) == "" {
			errs = append(errs, "file store S3 bucket name is required when S3 endpoint is set")
		}
		// In production require credentials
		if strings.EqualFold(mode, "production") {
			if config.FileStore.S3AccessKey == "" || config.FileStore.S3SecretKey == "" {
				errs = append(errs, "S3 access key and secret are required in production when S3 endpoint is set")
			}
		}
	}

	// Logging format & request body logging
	logFormat := strings.ToLower(strings.TrimSpace(config.Logging.Format))
	if logFormat != "" {
		validFormats := []string{"json", "pretty"}
		if !slices.Contains(validFormats, logFormat) {
			errs = append(errs, fmt.Sprintf("invalid log format: %q (valid: %v)", config.Logging.Format, validFormats))
		}
	}

	// Rate limiting
	if config.App.RateLimitEnabled {
		if config.App.RateLimitRequests <= 0 {
			errs = append(errs, "rate limit requests must be > 0 when rate limiting is enabled")
		}
		if config.App.RateLimitDuration <= 0 {
			errs = append(errs, "rate limit duration must be > 0 when rate limiting is enabled")
		}
	}

	// Telemetry
	if config.OTel.EnableTelemetry {
		if strings.TrimSpace(config.OTel.ExporterOTLPEndpoint) == "" {
			errs = append(errs, "OTel exporter endpoint is required when telemetry is enabled")
		}
		if config.OTel.TracingSampleRate < 0 || config.OTel.TracingSampleRate > 1 {
			errs = append(errs, fmt.Sprintf("OTel tracing sample rate must be between 0 and 1 (got %v)", config.OTel.TracingSampleRate))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}
