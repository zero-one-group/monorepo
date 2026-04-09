//! Config struct definitions — byte-for-byte port of
//! `apps/go-modular/internal/config/types.go` except for the two
//! deletions locked in the Phase D plan:
//!
//! - [`AppConfig`] drops `JWTAlgorithm` (design 3.4: HS256 only).
//! - [`MailerConfig`] drops `SMTPSecure` (design 3.5: lettre handles
//!   TLS mode via transport builder, not a bool config).
//!
//! Total fields: 12 + 3 + 6 + 8 + 3 + 7 = **39** across 6 sections.

use serde::{Deserialize, Serialize};

/// Top-level aggregate. Figment extracts this directly from env vars.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    #[serde(flatten)]
    pub app: AppConfig,
    #[serde(flatten)]
    pub database: DatabaseConfig,
    #[serde(flatten)]
    pub mailer: MailerConfig,
    #[serde(flatten)]
    pub file_store: FileStoreConfig,
    #[serde(flatten)]
    pub logging: LoggingConfig,
    #[serde(flatten)]
    pub otel: OTelConfig,
}

/// `AppConfig` — 12 fields (was 13 in Go; `jwt_algorithm` deleted).
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub app_mode: String,
    pub app_base_url: String,
    pub jwt_secret_key: String,
    pub server_host: String,
    pub server_port: u16,
    pub cors_origins: Vec<String>,
    pub cors_max_age: u32,
    pub cors_credentials: bool,
    pub rate_limit_enabled: bool,
    pub rate_limit_requests: u32,
    pub rate_limit_burst_size: u32,
    pub enable_api_docs: bool,
}

impl AppConfig {
    /// Accessor alias matching the Go source's `cfg.app.Mode` shape.
    #[inline]
    pub fn mode(&self) -> &str {
        &self.app_mode
    }

    #[inline]
    pub fn is_production(&self) -> bool {
        self.app_mode.eq_ignore_ascii_case("production")
    }
}

/// `DatabaseConfig` — 3 fields.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DatabaseConfig {
    pub database_url: String,
    pub pg_max_pool_size: u32,
    pub pg_max_retries: i32,
}

/// `MailerConfig` — 6 fields (was 7 in Go; `smtp_secure` deleted).
///
/// TLS mode is inferred from `smtp_port` in the lettre transport
/// builder (D-SMTP-1):
/// - 465 → implicit TLS (`Tls::Wrapper`)
/// - 587 → STARTTLS (required in production mode)
/// - 25, 1025 → plaintext (rejected in production by validator)
///
/// The `smtp_` prefix on every field is deliberate — env var names
/// require it (`SMTP_HOST`, `SMTP_PORT`, etc.) and figment maps struct
/// field names 1:1 to lowercased env keys.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(clippy::struct_field_names)]
pub struct MailerConfig {
    pub smtp_host: String,
    pub smtp_port: u16,
    pub smtp_username: String,
    pub smtp_password: String,
    pub smtp_sender_name: String,
    pub smtp_sender_email: String,
}

/// `FileStoreConfig` — 8 fields. Ported as dead config for forward
/// compatibility; no S3 client is instantiated by the auth/user paths
/// in this phase.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FileStoreConfig {
    pub public_assets_url: String,
    pub s3_endpoint: String,
    pub s3_access_key: String,
    pub s3_secret_key: String,
    pub s3_bucket_name: String,
    pub s3_region: String,
    pub s3_force_path_style: bool,
    pub s3_use_ssl: bool,
}

/// `LoggingConfig` — 3 fields. `log_` prefix on every field matches
/// the env var naming (`LOG_LEVEL`, `LOG_FORMAT`, `LOG_NO_COLOR`).
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(clippy::struct_field_names)]
pub struct LoggingConfig {
    pub log_level: String,
    pub log_format: String,
    pub log_no_color: bool,
}

/// `OTelConfig` — 7 fields. `otel_` prefix on every field matches the
/// env var naming (`OTEL_SERVICE_NAME`, `OTEL_EXPORTER_OTLP_*`, etc.).
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(clippy::struct_field_names)]
pub struct OTelConfig {
    pub otel_service_name: String,
    pub otel_exporter_otlp_protocol: String,
    pub otel_exporter_otlp_endpoint: String,
    pub otel_exporter_otlp_headers: String,
    pub otel_enable_telemetry: bool,
    pub otel_insecure_mode: bool,
    pub otel_tracing_sample_rate: f64,
}
