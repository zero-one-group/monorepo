//! Default config values — byte-for-byte port of
//! `apps/{{ package_name | kebab_case }}/internal/config/default.go`, except for the two
//! deletions (see `types.rs`).

use super::types::{
    AppConfig, Config, DatabaseConfig, FileStoreConfig, LoggingConfig, MailerConfig, OTelConfig,
};

pub fn default_config() -> Config {
    Config {
        app: AppConfig {
            app_mode: "development".to_string(),
            app_base_url: "http://localhost:{{ port_number }}".to_string(),
            jwt_secret_key: "_THIS_IS_DEFAULT_JWT_SECRET_KEY_".to_string(),
            server_host: "0.0.0.0".to_string(),
            server_port: {{ port_number }},
            cors_origins: vec!["*".to_string()],
            cors_max_age: 300,
            cors_credentials: true,
            rate_limit_enabled: true,
            rate_limit_requests: 20,
            rate_limit_burst_size: 60,
            enable_api_docs: true,
        },
        database: DatabaseConfig {
            database_url: "postgresql://postgres:securedb@localhost:5432/postgres?sslmode=disable"
                .to_string(),
            pg_max_pool_size: 10,
            pg_max_retries: 5,
        },
        mailer: MailerConfig {
            smtp_host: "localhost".to_string(),
            smtp_port: 1025,
            smtp_username: String::new(),
            smtp_password: String::new(),
            // Note: Go defaults had escaped quotes ("\"System Mailer\"").
            // We strip them here per audit §9.12 — viper/godotenv strip
            // quotes on load, so the effective Go value was also plain
            // "System Mailer" at runtime.
            smtp_sender_name: "System Mailer".to_string(),
            smtp_sender_email: "mailer@example.com".to_string(),
        },
        file_store: FileStoreConfig {
            public_assets_url: "http://localhost:8010".to_string(),
            s3_endpoint: "http://localhost:9100".to_string(),
            s3_access_key: "s3admin".to_string(),
            s3_secret_key: "s3passw0rd".to_string(),
            s3_bucket_name: "devbucket".to_string(),
            s3_region: "auto".to_string(),
            s3_force_path_style: false,
            s3_use_ssl: false,
        },
        logging: LoggingConfig {
            log_level: "info".to_string(),
            log_format: "pretty".to_string(),
            log_no_color: false,
        },
        otel: OTelConfig {
            otel_service_name: "{{ package_name | kebab_case }}".to_string(),
            otel_exporter_otlp_protocol: "http/protobuf".to_string(),
            otel_exporter_otlp_endpoint: "http://localhost:4318".to_string(),
            otel_exporter_otlp_headers: "authorization=YOUR_INGESTION_API_KEY".to_string(),
            otel_enable_telemetry: false,
            otel_insecure_mode: false,
            otel_tracing_sample_rate: 0.7,
        },
    }
}
