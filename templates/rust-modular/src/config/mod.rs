//! Config loader for {{ package_name | kebab_case }}.
//!
//! Ports `apps/{{ package_name | kebab_case }}/internal/config/` from Go to Rust:
//!
//! - `types.rs` mirrors `types.go` (6 sections, 39 fields after deleting
//!   `JWTAlgorithm` and `SMTPSecure` per plan design decisions 3.4 and 3.5).
//! - `defaults.rs` mirrors `default.go`.
//! - `validate.rs` mirrors `validate.go` plus Phase D extras (reject
//!   plaintext SMTP in production mode).
//!
//! Env var names are 1:1 with the Go source. Figment auto-lowercases
//! env keys (Phase B finding), so the Rust structs use idiomatic
//! `snake_case` field names with no `#[serde(rename)]`.

pub mod defaults;
pub mod types;
pub mod validate;

pub use types::{
    AppConfig, Config, DatabaseConfig, FileStoreConfig, LoggingConfig, MailerConfig, OTelConfig,
};

use anyhow::{Context, Result};
use figment::Figment;
use figment::providers::{Env as FigEnv, Serialized};

impl Config {
    /// Load the full config from env vars, layered over static defaults.
    ///
    /// Ordering: defaults (from [`defaults::default_config`]) are merged
    /// first, then env vars override via `Figment::merge(Env::raw())`.
    /// Unknown env keys are ignored (matches viper behaviour in Go).
    pub fn from_environment() -> Result<Self> {
        let figment = Figment::new()
            .merge(Serialized::defaults(defaults::default_config()))
            .merge(FigEnv::raw());
        let config: Self = figment.extract().context("extract config from figment")?;
        validate::validate_config(&config).context("validate config")?;
        Ok(config)
    }

    /// Test-only helper: start from defaults, validate, return.
    pub fn from_defaults() -> Result<Self> {
        let config = defaults::default_config();
        validate::validate_config(&config).context("validate default config")?;
        Ok(config)
    }

    /// `true` when `App.mode == "production"` (case-insensitive).
    pub fn is_production(&self) -> bool {
        self.app.is_production()
    }
}
