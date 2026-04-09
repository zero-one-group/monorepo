//! Env loader. Mirrors the env var surface of
//! `apps/go-clean/.env.example` + `apps/go-clean/config/`.
//!
//! Figment auto-lowercases env keys (per the Phase B finding in
//! `apps/fastapi-ai/src/core/env.rs`), so the Rust struct uses
//! idiomatic `snake_case` field names and no `#[serde(rename)]`.

use anyhow::Result;
use figment::Figment;
use figment::providers::Env as FigEnv;
use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct Env {
    #[serde(default = "default_service_name")]
    pub service_name: String,

    #[serde(default = "default_app_environment")]
    pub app_environment: String,

    #[serde(default = "default_app_host")]
    pub app_host: String,

    #[serde(default = "default_app_port")]
    pub app_port: u16,

    pub database_url: String,

    #[serde(default)]
    pub cors_allow_origins: Option<String>,

    pub jwt_secret: String,

    /// Access token TTL in minutes. Matches the env var name the Go
    /// `utils/jwt.go` actually reads (`AUTH_TOKEN_EXPIRY_MINUTES`),
    /// NOT the mis-documented `JWT_TOKEN_EXPIRY_MINUTES` in
    /// `.env.example`. Go's code falls back to 60 if unset/invalid,
    /// so we preserve that.
    #[serde(default = "default_auth_token_expiry_minutes")]
    pub auth_token_expiry_minutes: i64,

    #[serde(default)]
    pub enable_swagger: bool,

    #[serde(default = "default_otel_endpoint")]
    pub otel_exporter_otlp_endpoint: String,
}

fn default_service_name() -> String {
    "go-clean".to_string()
}
fn default_app_environment() -> String {
    "local".to_string()
}
fn default_app_host() -> String {
    "127.0.0.1".to_string()
}
fn default_app_port() -> u16 {
    8000
}
fn default_auth_token_expiry_minutes() -> i64 {
    60
}
fn default_otel_endpoint() -> String {
    "localhost:4317".to_string()
}

impl Env {
    pub fn from_environment() -> Result<Self> {
        let figment = Figment::new().merge(FigEnv::raw());
        let env: Self = figment.extract()?;
        Ok(env)
    }

    pub fn is_production(&self) -> bool {
        self.app_environment == "production"
    }
}
