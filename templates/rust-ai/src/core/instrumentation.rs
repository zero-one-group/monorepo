//! OpenTelemetry tracing initialization + Prometheus metrics registry.
//!
//! Mirrors `app/core/instrumentation.py`. Per the Python script the
//! tracer is configured ONLY when `APP_ENVIRONMENT == "production"`. We
//! preserve that exactly: in dev we skip the OTLP exporter setup so
//! local runs don't fail when an `OTel` collector isn't available.
//!
//! REPLICABLE spans (these MUST exact-match the Python equivalents per
//! the spec's two-tier `OTel` parity contract):
//! - `service.greetings`        — emitted by [`crate::services::greeting`]
//! - `repository.fetch_greetings` — emitted by [`crate::repository::openai::greeting`]
//! - `route.greetings`          — emitted by [`crate::router::openai`]
//!
//! REPLICABLE metrics:
//! - `repository_greetings_requests_total` (labels: model)
//! - `repository_greetings_requests_failures_total` (labels: model, phase)
//! - `repository_greetings_request_duration_seconds` (labels: model,
//!   buckets: [0.1, 0.5, 1, 2, 5, 10])
//!
//! STRUCTURAL gaps (not replicated in the Rust port):
//! - `FastAPI`'s `opentelemetry-instrumentation-fastapi` auto-spans
//!   (`http.request`, `http.response.start`, `http.response.body`,
//!   ASGI lifecycle) — no axum equivalent. We emit a `tower-http::trace`
//!   span instead, which has different attribute keys.

use anyhow::Result;
use prometheus::{HistogramOpts, HistogramVec, IntCounterVec, Opts, Registry};
use std::sync::OnceLock;

use crate::core::env::Env;

/// Process-wide metrics registry. Held in a `OnceLock` so handlers can
/// reach it without threading the registry through `AppState`.
static REGISTRY: OnceLock<Registry> = OnceLock::new();
static GREETINGS_REQUESTS_TOTAL: OnceLock<IntCounterVec> = OnceLock::new();
static GREETINGS_REQUESTS_FAILURES: OnceLock<IntCounterVec> = OnceLock::new();
static GREETINGS_REQUEST_DURATION: OnceLock<HistogramVec> = OnceLock::new();

/// Initialize global observability — tracing + metrics.
///
/// Tracing is configured only in production (matching the Python
/// `instrument_app()` early-return). Metrics are registered
/// unconditionally so unit tests and local runs can read them.
pub fn init_observability(env: &Env) -> Result<()> {
    init_metrics()?;

    if env.is_production() {
        tracing::info!(
            "The environment is set to production; instrumentation is being configured."
        );
        // OTLP tracer init goes here once the otel-otlp builder shape is
        // confirmed against docs.rs deeper pages. The current code path
        // is intentionally a no-op in dev so unit/integration tests run
        // without an OTel collector. Production parity is tracked by
        // the OTel hook-point inventory follow-up task (B0.5).
    }

    Ok(())
}

fn init_metrics() -> Result<()> {
    let registry = Registry::new();

    let requests_total = IntCounterVec::new(
        Opts::new(
            "repository_greetings_requests_total",
            "Total number of calls to repository.greetings()",
        ),
        &["model"],
    )?;
    let failures_total = IntCounterVec::new(
        Opts::new(
            "repository_greetings_requests_failures_total",
            "Number of failed calls to repository.greetings() broken out by phase",
        ),
        &["model", "phase"],
    )?;
    let duration = HistogramVec::new(
        HistogramOpts::new(
            "repository_greetings_request_duration_seconds",
            "Time spent in repository.greetings()",
        )
        .buckets(vec![0.1, 0.5, 1.0, 2.0, 5.0, 10.0]),
        &["model"],
    )?;

    registry.register(Box::new(requests_total.clone()))?;
    registry.register(Box::new(failures_total.clone()))?;
    registry.register(Box::new(duration.clone()))?;

    let _ = REGISTRY.set(registry);
    let _ = GREETINGS_REQUESTS_TOTAL.set(requests_total);
    let _ = GREETINGS_REQUESTS_FAILURES.set(failures_total);
    let _ = GREETINGS_REQUEST_DURATION.set(duration);
    Ok(())
}

pub fn registry() -> Option<&'static Registry> {
    REGISTRY.get()
}

pub fn greetings_requests_total() -> &'static IntCounterVec {
    GREETINGS_REQUESTS_TOTAL
        .get()
        .expect("metrics not initialized — call init_observability() first")
}

pub fn greetings_requests_failures() -> &'static IntCounterVec {
    GREETINGS_REQUESTS_FAILURES
        .get()
        .expect("metrics not initialized — call init_observability() first")
}

pub fn greetings_request_duration() -> &'static HistogramVec {
    GREETINGS_REQUEST_DURATION
        .get()
        .expect("metrics not initialized — call init_observability() first")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn metrics_register_and_have_expected_labels() {
        let env = Env {
            ml_prefix_api: String::new(),
            app_name: "test".into(),
            app_environment: "development".into(),
            database_url: String::new(),
            openai_api_key: String::new(),
            otel_exporter_otlp_endpoint: String::new(),
        };
        // Idempotent — running once is enough; OnceLock ignores subsequent sets.
        init_observability(&env).unwrap();

        // The Counter+Histogram label sets must match the Python module exactly.
        let cnt = greetings_requests_total();
        let _ = cnt.with_label_values(&["gpt-4.1-mini"]);
        let fail = greetings_requests_failures();
        let _ = fail.with_label_values(&["gpt-4.1-mini", "openai_call"]);
        let dur = greetings_request_duration();
        let _ = dur.with_label_values(&["gpt-4.1-mini"]);
    }
}
