//! Config validation — port of `internal/config/validate.go` plus two
//! Phase D extras:
//!
//! 1. Reject plaintext SMTP ports (25, 1025) in production mode.
//!    Complements the design 3.5 deletion of `SMTPSecure`: the
//!    lettre transport builder enforces TLS on 587/465; plaintext
//!    is only allowed in dev.
//!
//! 2. No RS256 / `JWT_ALGORITHM` validation (design 3.4: HS256 only).

use anyhow::{Result, bail};

use super::types::Config;

/// Validate the fully-loaded config. Returns a joined error with all
/// failing checks so operators see every problem in one boot attempt.
pub fn validate_config(config: &Config) -> Result<()> {
    let mut errs: Vec<String> = Vec::new();

    // App mode
    let mode = config.app.app_mode.trim().to_ascii_lowercase();
    if mode.is_empty() {
        errs.push("app mode is required".to_string());
    } else if !matches!(
        mode.as_str(),
        "development" | "production" | "staging" | "test"
    ) {
        errs.push(format!(
            "invalid app mode: \"{}\" (valid: development, production, staging, test)",
            config.app.app_mode
        ));
    }

    // JWT secret
    let secret = config.app.jwt_secret_key.trim();
    if mode == "production" && (secret.is_empty() || secret == "_THIS_IS_DEFAULT_JWT_SECRET_KEY_") {
        errs.push("JWT secret key must be set in production".to_string());
    }

    // Server port
    if config.app.server_port == 0 {
        errs.push(format!(
            "invalid server port: {} (must be 1-65535)",
            config.app.server_port
        ));
    }

    // Database URL — loose check. sqlx will return a clearer error on
    // actual connect; validator only catches empty + obviously wrong
    // scheme. Matches Go's `url.Parse` which was very permissive.
    let db_url = config.database.database_url.trim();
    if db_url.is_empty() {
        errs.push("database URL is required".to_string());
    } else if !(db_url.starts_with("postgres://") || db_url.starts_with("postgresql://")) {
        errs.push(format!(
            "invalid database URL: \"{db_url}\" (must start with postgres:// or postgresql://)"
        ));
    }

    // Postgres pool
    if config.database.pg_max_pool_size == 0 {
        errs.push(format!(
            "pg max pool size must be > 0 (got {})",
            config.database.pg_max_pool_size
        ));
    }
    if config.database.pg_max_retries < -1 {
        errs.push(format!(
            "pg max retries must be >= -1 (got {}); -1 means infinite retry",
            config.database.pg_max_retries
        ));
    }

    // Mailer host/port
    let smtp_host = config.mailer.smtp_host.trim();
    if !smtp_host.is_empty() && config.mailer.smtp_port == 0 {
        errs.push(format!(
            "invalid mailer SMTP port: {} (must be 1-65535)",
            config.mailer.smtp_port
        ));
    }

    // Phase D extra: reject non-TLS SMTP ports in production.
    // Ports 587 (STARTTLS) and 465 (implicit TLS) are the only allowed
    // ports in production mode. Anything else — including the common
    // 25/1025 dev plaintext ports — is rejected to surface
    // misconfiguration at boot rather than silently sending plaintext.
    if mode == "production" && !smtp_host.is_empty() {
        let port = config.mailer.smtp_port;
        if port != 587 && port != 465 {
            errs.push(format!(
                "SMTP port {port} forbidden in production; use 587 (STARTTLS) or 465 (implicit TLS)"
            ));
        }
    }

    // FileStore
    let s3_ep = config.file_store.s3_endpoint.trim();
    if !s3_ep.is_empty() {
        if config.file_store.s3_bucket_name.trim().is_empty() {
            errs.push("file store S3 bucket name is required when S3 endpoint is set".to_string());
        }
        if mode == "production"
            && (config.file_store.s3_access_key.is_empty()
                || config.file_store.s3_secret_key.is_empty())
        {
            errs.push(
                "S3 access key and secret are required in production when S3 endpoint is set"
                    .to_string(),
            );
        }
    }

    // Logging format
    let log_format = config.logging.log_format.trim().to_ascii_lowercase();
    if !log_format.is_empty() && !matches!(log_format.as_str(), "json" | "pretty") {
        errs.push(format!(
            "invalid log format: \"{}\" (valid: json, pretty)",
            config.logging.log_format
        ));
    }

    // Rate limiting
    if config.app.rate_limit_enabled {
        if config.app.rate_limit_requests == 0 {
            errs.push("rate limit requests must be > 0 when rate limiting is enabled".to_string());
        }
        if config.app.rate_limit_burst_size == 0 {
            errs.push(
                "rate limit burst size must be > 0 when rate limiting is enabled".to_string(),
            );
        }
    }

    // OTel
    if config.otel.otel_enable_telemetry {
        if config.otel.otel_exporter_otlp_endpoint.trim().is_empty() {
            errs.push("OTel exporter endpoint is required when telemetry is enabled".to_string());
        }
        if !(0.0..=1.0).contains(&config.otel.otel_tracing_sample_rate) {
            errs.push(format!(
                "OTel tracing sample rate must be between 0 and 1 (got {})",
                config.otel.otel_tracing_sample_rate
            ));
        }
    }

    if !errs.is_empty() {
        bail!("config validation failed: {}", errs.join("; "));
    }
    Ok(())
}
