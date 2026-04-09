//! Environment loading.
//!
//! Mirrors `app/core/env.py` field-for-field. Loaded once at startup
//! via figment + dotenvy and threaded through `AppState`.
//!
//! **Figment env key convention**: `figment::providers::Env::raw()`
//! automatically lowercases every environment variable name before
//! deserializing, so `ML_PREFIX_API` (the shell env key) becomes
//! `ml_prefix_api` (the struct field). We rely on this default so the
//! Rust struct can use idiomatic `snake_case` field names without any
//! `#[serde(rename = "...")]` noise. The Python service uses
//! `UPPER_SNAKE` as its env var names (via pydantic-settings) and the
//! shell/`.env` files use `UPPER_SNAKE` too — figment's lowercase
//! transform bridges those conventions transparently.

use anyhow::Result;
use figment::Figment;
use figment::providers::Env as FigEnv;
use serde::Deserialize;

/// Application settings, populated from env vars.
///
/// Env var mapping (figment lowercases the shell key automatically):
/// - `ML_PREFIX_API` → `ml_prefix_api`
/// - `APP_NAME` → `app_name` (default `"fastapi-ai"`)
/// - `APP_ENVIRONMENT` → `app_environment` (default `"development"`)
/// - `DATABASE_URL` → `database_url`
/// - `OPENAI_API_KEY` → `openai_api_key`
/// - `OTEL_EXPORTER_OTLP_ENDPOINT` → `otel_exporter_otlp_endpoint`
#[derive(Debug, Clone, Deserialize)]
pub struct Env {
    pub ml_prefix_api: String,

    #[serde(default = "default_app_name")]
    pub app_name: String,

    #[serde(default = "default_app_environment")]
    pub app_environment: String,

    pub database_url: String,
    pub openai_api_key: String,
    pub otel_exporter_otlp_endpoint: String,
}

fn default_app_name() -> String {
    "fastapi-ai".to_string()
}

fn default_app_environment() -> String {
    "development".to_string()
}

impl Env {
    /// Load from process environment.
    pub fn from_environment() -> Result<Self> {
        let figment = Figment::new().merge(FigEnv::raw());
        let env: Self = figment.extract()?;
        Ok(env)
    }

    /// Convenience: are we running in production?
    pub fn is_production(&self) -> bool {
        self.app_environment == "production"
    }
}
