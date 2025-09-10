package config

// JWTAlgorithm is a typesafe enum for JWT algorithm
// Supported values: "HS256", "RS256"
type JWTAlgorithm string

const (
	JWTAlgorithmHS256 JWTAlgorithm = "HS256"
	JWTAlgorithmRS256 JWTAlgorithm = "RS256"
)

type Config struct {
	App       AppConfig       `env:",squash"`
	Database  DatabaseConfig  `env:",squash"`
	Mailer    MailerConfig    `env:",squash"`
	FileStore FileStoreConfig `env:",squash"`
	Logging   LoggingConfig   `env:",squash"`
	OTel      OTelConfig      `env:",squash"`
}

type AppConfig struct {
	Mode              string       `env:"APP_MODE"` // development|production
	BaseURL           string       `env:"APP_BASE_URL"`
	JWTSecretKey      string       `env:"JWT_SECRET_KEY"`
	JWTAlgorithm      JWTAlgorithm `env:"JWT_ALGORITHM"`
	ServerHost        string       `env:"SERVER_HOST"`
	ServerPort        int          `env:"SERVER_PORT"`
	CORSOrigins       []string     `env:"CORS_ORIGINS"`
	CORSMaxAge        int          `env:"CORS_MAX_AGE"`
	CORSCredentials   bool         `env:"CORS_CREDENTIALS"`
	RateLimitEnabled  bool         `env:"RATE_LIMIT_ENABLED"`
	RateLimitRequests int          `env:"RATE_LIMIT_REQUESTS"`
	RateLimitDuration int          `env:"RATE_LIMIT_DURATION"`
	EnableAPIDocs     bool         `env:"ENABLE_API_DOCS"`
}

type DatabaseConfig struct {
	PostgresURL   string `env:"DATABASE_URL"`
	PgMaxPoolSize int    `env:"PG_MAX_POOL_SIZE"`
	PgMaxRetries  int    `env:"PG_MAX_RETRIES"`
}

type MailerConfig struct {
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SenderName   string `env:"SMTP_SENDER_NAME"`
	SenderEmail  string `env:"SMTP_SENDER_EMAIL"`
	SMTPSecure   bool   `env:"SMTP_SECURE"`
}

type FileStoreConfig struct {
	PublicAssetsURL  string `env:"PUBLIC_ASSETS_URL"`
	S3Endpoint       string `env:"S3_ENDPOINT"`
	S3AccessKey      string `env:"S3_ACCESS_KEY"`
	S3SecretKey      string `env:"S3_SECRET_KEY"`
	S3BucketName     string `env:"S3_BUCKET_NAME"`
	S3Region         string `env:"S3_REGION"`
	S3ForcePathStyle bool   `env:"S3_FORCE_PATH_STYLE"`
	S3UseSSL         bool   `env:"S3_USE_SSL"`
}

type LoggingConfig struct {
	Level   string `env:"LOG_LEVEL"`
	Format  string `env:"LOG_FORMAT"`
	NoColor bool   `env:"LOG_NO_COLOR"`
}

type OTelConfig struct {
	ServiceName          string  `env:"OTEL_SERVICE_NAME"`
	ExporterOTLPProtocol string  `env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
	ExporterOTLPEndpoint string  `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	ExporterOTLPHeaders  string  `env:"OTEL_EXPORTER_OTLP_HEADERS"`
	EnableTelemetry      bool    `env:"OTEL_ENABLE_TELEMETRY"`
	InsecureMode         bool    `env:"OTEL_INSECURE_MODE"`
	TracingSampleRate    float64 `env:"OTEL_TRACING_SAMPLE_RATE"`
}
