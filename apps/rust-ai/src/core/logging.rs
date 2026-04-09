//! Logging initialization and per-request `RequestId` middleware.
//!
//! The Python service uses `python-json-logger` with JSON output in
//! production and plain text in development. We do the same with
//! `tracing-subscriber`'s `fmt` layer (json feature in production,
//! plain in dev).
//!
//! Request IDs are wired via `tower-http::request_id::SetRequestIdLayer`
//! plus `PropagateRequestIdLayer`, both layered onto the axum router.

use tracing_subscriber::EnvFilter;
use tracing_subscriber::fmt::format::FmtSpan;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;

use crate::core::env::Env;

/// Initialize global logging.
///
/// Idempotent — calling more than once is a no-op (the underlying
/// `tracing_subscriber` registry rejects duplicate global subscribers,
/// so we ignore the error).
pub fn init_logging(env: &Env) {
    let filter = EnvFilter::try_from_default_env()
        .unwrap_or_else(|_| EnvFilter::new(if env.is_production() { "info" } else { "debug" }));

    let registry = tracing_subscriber::registry().with(filter);

    if env.is_production() {
        let _ = registry
            .with(
                tracing_subscriber::fmt::layer()
                    .json()
                    .with_span_events(FmtSpan::CLOSE),
            )
            .try_init();
    } else {
        let _ = registry
            .with(tracing_subscriber::fmt::layer().with_target(true))
            .try_init();
        tracing::debug!("Debug logging enabled");
    }
}
