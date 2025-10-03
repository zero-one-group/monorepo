package config

func DefaultConfig() Config {
	return Config{
		App: AppConfig{
			Mode:               "development",
			BaseURL:            "http://localhost:{{ port_number }}",
			JWTSecretKey:       "_THIS_IS_DEFAULT_JWT_SECRET_KEY_",
			JWTAlgorithm:       JWTAlgorithmHS256,
			ServerHost:         "0.0.0.0",
			ServerPort:         {{ port_number }},
			CORSOrigins:        []string{"*"},
			CORSMaxAge:         300,
			CORSCredentials:    true,
			RateLimitEnabled:   true,
			RateLimitRequests:  20,
			RateLimitBurstSize: 60,
			EnableAPIDocs:      true,
		},
		Database: DatabaseConfig{
			PostgresURL:   "postgresql://postgres:securedb@localhost:5432/postgres?sslmode=disable",
			PgMaxPoolSize: 10,
			PgMaxRetries:  5,
		},
		Mailer: MailerConfig{
			SMTPHost:     "localhost",
			SMTPPort:     1025,
			SMTPUsername: "",
			SMTPPassword: "",
			SenderName:   "\"System Mailer\"",
			SenderEmail:  "\"mailer@example.com\"",
			SMTPSecure:   false,
		},
		FileStore: FileStoreConfig{
			PublicAssetsURL:  "http://localhost:8010",
			S3Endpoint:       "http://localhost:9100",
			S3Region:         "auto",
			S3AccessKey:      "s3admin",
			S3SecretKey:      "s3passw0rd",
			S3BucketName:     "devbucket",
			S3ForcePathStyle: false,
			S3UseSSL:         false,
		},
		Logging: LoggingConfig{
			Level:   "info",
			Format:  "pretty",
			NoColor: false,
		},
		OTel: OTelConfig{
			ExporterOTLPProtocol: "http/protobuf",
			ExporterOTLPEndpoint: "http://localhost:4318",
			ExporterOTLPHeaders:  "\"authorization=YOUR_INGESTION_API_KEY\"",
			ServiceName:          "{{ package_name | kebab_case }}",
			EnableTelemetry:      false,
			InsecureMode:         false,
			TracingSampleRate:    0.7,
		},
	}
}
