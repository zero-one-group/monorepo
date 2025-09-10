package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strings"

	"github.com/labstack/gommon/log"
)

var (
	rePasswordQuery = regexp.MustCompile(`(?i)(password|pwd)=([^&\s]+)`)
	reUserPassAt    = regexp.MustCompile(`(?i)(://)?([^:@/]+):([^@/]+)@`)
	maskPlaceholder = "*****"
)

// Returns true if running in production environment
func (c *Config) IsProduction() bool {
	if c == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(c.App.Mode), "production")
}

// Returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	if c == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(c.App.Mode), "development")
}

// Returns true if debug mode is enabled (based on log level)
func (c *Config) IsDebug() bool {
	if c == nil {
		return false
	}
	logLevel := strings.ToLower(strings.TrimSpace(c.Logging.Level))
	return strings.Contains(logLevel, "debug") || strings.Contains(logLevel, "trace")
}

func (c *Config) GetAppBaseURL() string {
	if c == nil {
		return ""
	}
	return c.App.BaseURL
}

// Returns the http server address in host:port format
func (c *Config) GetServerAddr() string {
	if c == nil {
		return ""
	}
	host := strings.TrimSpace(c.App.ServerHost)
	if host == "" {
		// If host is empty, return :port to bind on all interfaces
		return fmt.Sprintf(":%d", c.App.ServerPort)
	}
	return fmt.Sprintf("%s:%d", host, c.App.ServerPort)
}

// Returns true if API Docs are enabled
func (c *Config) IsAPIDocsEnabled() bool {
	if c == nil {
		return false
	}
	return c.App.EnableAPIDocs
}

// Returns true if telemetry is enabled
func (c *Config) IsTelemetryEnabled() bool {
	if c == nil {
		return false
	}
	return c.OTel.EnableTelemetry
}

// GetTelemetryStatus returns "enabled" or "disabled" based on telemetry configuration
func (c *Config) GetTelemetryStatus() string {
	if c == nil {
		return "disabled"
	}
	if c.IsTelemetryEnabled() {
		return "enabled"
	}
	return "disabled"
}

// Returns the database URL in the format expected by the database driver
func (c *Config) GetDatabaseURL() string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.Database.PostgresURL)
}

// GetMaskedDatabaseURL returns the DB URL with password masked.
// Handles both URL-style (scheme://user:pass@host...) and query-style (password=...).
func (c *Config) GetMaskedDatabaseURL() string {
	rawURL := c.GetDatabaseURL()
	if rawURL == "" {
		return ""
	}

	// Try URL parse first
	u, err := url.Parse(rawURL)
	if err == nil && u.User != nil {
		username := u.User.Username()
		if _, hasPassword := u.User.Password(); hasPassword {
			u.User = url.UserPassword(username, maskPlaceholder)
			masked := u.String()
			// decode url-encoded asterisks if any
			masked = strings.ReplaceAll(masked, "%2A", "*")
			return masked
		}
		return u.String()
	}

	// Fallback: mask query-style password or pwd params
	if rePasswordQuery.MatchString(rawURL) {
		return rePasswordQuery.ReplaceAllString(rawURL, "${1}="+maskPlaceholder)
	}

	// Fallback: mask user:pass@ patterns even if url.Parse failed
	if reUserPassAt.MatchString(rawURL) {
		return reUserPassAt.ReplaceAllString(rawURL, "${1}${2}:"+maskPlaceholder+"@")
	}

	// If nothing matched, return original
	return rawURL
}

// MaskJWTSecret returns JWT secret with most characters masked.
// Reveals a small prefix/suffix for easier debugging but hides the rest.
func (c *Config) MaskJWTSecret() string {
	if c == nil {
		return "***"
	}
	secret := strings.TrimSpace(c.App.JWTSecretKey)
	if secret == "" {
		return ""
	}
	if len(secret) <= 8 {
		return maskPlaceholder
	}
	// reveal first 4 and last 4
	return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}

// GetLogLevel convert to slog.Level
func (c *Config) GetLogLevel() slog.Leveler {
	if c == nil {
		return slog.LevelInfo
	}
	logLevel := strings.ToLower(strings.TrimSpace(c.Logging.Level))
	switch logLevel {
	case "debug", "trace":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// GetEchoLogLevel converts configured log level to github.com/labstack/gommon/log.Lvl
// Aligned with GetLogLevel: no single-letter shortcuts, map trace->debug, default to INFO.
func (c *Config) GetEchoLogLevel() log.Lvl {
	if c == nil {
		return log.INFO
	}
	level := strings.ToLower(strings.TrimSpace(c.Logging.Level))

	switch level {
	case "debug", "trace":
		return log.DEBUG
	case "info", "":
		return log.INFO
	case "warn", "warning":
		return log.WARN
	case "error":
		return log.ERROR
	default:
		return log.INFO
	}
}

// GetLogFormat returns the configured log format
func (c *Config) GetLogFormat() string {
	if c == nil {
		return "pretty"
	}
	return strings.ToLower(strings.TrimSpace(c.Logging.Format))
}

// IsColorDisabled returns true if colors should be disabled
func (c *Config) IsColorDisabled() bool {
	if c == nil {
		return false
	}
	return c.Logging.NoColor
}
