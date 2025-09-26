package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfigAndGetters(t *testing.T) {
	t.Run("DefaultConfig_is_valid", func(t *testing.T) {
		cfg := DefaultConfig()
		err := validateConfig(&cfg)
		require.NoError(t, err)
	})

	t.Run("InvalidAppMode", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.App.Mode = "bad-mode"
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid app mode")
	})

	t.Run("Production_requires_jwt_secret", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.App.Mode = "production"
		// Default secret is the placeholder and should be rejected in production
		cfg.App.JWTSecretKey = "_THIS_IS_DEFAULT_JWT_SECRET_KEY_"
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "JWT secret key must be set in production")
	})

	t.Run("InvalidDatabaseURL", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Database.PostgresURL = "://not-a-valid-url"
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid database URL")
	})

	t.Run("PgPoolSize_and_Retries_validation", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Database.PgMaxPoolSize = 0
		cfg.Database.PgMaxRetries = -1
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pg max pool size must be > 0")
		assert.Contains(t, err.Error(), "pg max retries must be >= 0")
	})

	t.Run("Mailer_port_without_host", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Mailer.SMTPHost = ""
		cfg.Mailer.SMTPPort = 2525
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mailer SMTP host is empty but SMTP port is set")
	})

	t.Run("S3_endpoint_requires_bucket_and_creds_in_prod", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.FileStore.S3Endpoint = "http://s3.local"
		cfg.FileStore.S3BucketName = ""
		cfg.App.Mode = "development"
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file store S3 bucket name is required")

		// In production also require credentials
		cfg.FileStore.S3BucketName = "bucket"
		cfg.App.Mode = "production"
		cfg.FileStore.S3AccessKey = ""
		cfg.FileStore.S3SecretKey = ""
		err = validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "S3 access key and secret are required in production")
	})

	t.Run("Telemetry_settings_validation", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.OTel.EnableTelemetry = true
		cfg.OTel.ExporterOTLPEndpoint = ""
		cfg.OTel.TracingSampleRate = 1.5
		err := validateConfig(&cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "OTel exporter endpoint is required")
		assert.Contains(t, err.Error(), "OTel tracing sample rate must be between 0 and 1")
	})
}

func TestConfigHelpers(t *testing.T) {
	t.Run("GetServerAddr_and_modes", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.App.ServerPort = 1234
		cfg.App.ServerHost = ""
		addr := cfg.GetServerAddr()
		assert.Equal(t, ":1234", addr)

		cfg.App.Mode = "production"
		assert.True(t, cfg.IsProduction())
		assert.False(t, cfg.IsDevelopment())
	})

	t.Run("MaskDatabaseURL_variants", func(t *testing.T) {
		cfg := DefaultConfig()

		// URL-style with user:pass@
		cfg.Database.PostgresURL = "postgresql://alice:secretpass@localhost:5432/db"
		masked := cfg.GetMaskedDatabaseURL()
		assert.Contains(t, masked, "*****")
		assert.NotContains(t, masked, "secretpass")

		// Query-style password
		cfg.Database.PostgresURL = "http://example.local/path?password=verysecret&x=1"
		masked = cfg.GetMaskedDatabaseURL()
		assert.Contains(t, masked, "password=*****")
		assert.NotContains(t, masked, "verysecret")

		// Fallback user:pass@ pattern (non-URL)
		cfg.Database.PostgresURL = "user:superpass@somehost"
		masked = cfg.GetMaskedDatabaseURL()
		assert.Contains(t, masked, ":*****@")
	})

	t.Run("MaskJWTSecret_behavior", func(t *testing.T) {
		cfg := DefaultConfig()
		// short secret
		cfg.App.JWTSecretKey = "short"
		assert.Equal(t, maskPlaceholder, cfg.MaskJWTSecret())

		// longer secret
		cfg.App.JWTSecretKey = "abcdefghijkl"
		// expect first4 + stars + last4
		got := cfg.MaskJWTSecret()
		assert.Contains(t, got, "abcd")
		assert.Contains(t, got, "ijkl")
		assert.NotContains(t, got, "abcdefghij") // original secret not fully present
	})

	t.Run("Logging_and_telemetry_getters", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Logging.Level = "debug"
		cfg.Logging.Format = "json"
		cfg.Logging.NoColor = true

		assert.Equal(t, "json", cfg.GetLogFormat())
		assert.True(t, cfg.IsColorDisabled())
		// slog level and echo conversion should not panic
		assert.NotNil(t, cfg.GetLogLevel())
		assert.NotNil(t, cfg.GetSlogLevel())
		assert.NotNil(t, cfg.GetEchoLogLevel())

		// telemetry helpers
		cfg.OTel.EnableTelemetry = false
		assert.False(t, cfg.IsTelemetryEnabled())
		assert.Equal(t, "disabled", cfg.GetTelemetryStatus())
		cfg.OTel.EnableTelemetry = true
		assert.Equal(t, "enabled", cfg.GetTelemetryStatus())
	})
}
