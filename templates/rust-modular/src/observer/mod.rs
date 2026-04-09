//! Observability: tracing + `OTel` bootstrap.
//!
//! Port of `internal/observer/logger/logger.go` + `tracer.go` + `metrics.go`.
//!
//! This scaffold wires `tracing_subscriber` with env-filter from
//! `LOG_LEVEL`. Full `OTel` OTLP exporter wiring (gated on
//! `otel_enable_telemetry`) is stubbed — it follows the same pattern
//! as Phase B fastapi-ai and will be fleshed out during D-AUTH work
//! when the real span coverage is ready. For now the scaffold just
//! initializes tracing so middleware spans have somewhere to go.

use tracing::info;
use tracing_subscriber::EnvFilter;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;

use crate::config::{LoggingConfig, OTelConfig};

/// Initialize tracing subscriber (idempotent — safe to call multiple
/// times under tests; second call is a no-op via `try_init`).
pub fn init_tracing(logging: &LoggingConfig, otel: &OTelConfig) {
    let default_level = logging.log_level.as_str();
    let filter =
        EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new(default_level));

    let fmt_layer = tracing_subscriber::fmt::layer()
        .with_target(true)
        .with_ansi(!logging.log_no_color);

    // JSON or pretty output based on LOG_FORMAT. Pretty is the default
    // for local dev; json is for production log aggregators.
    let registry = tracing_subscriber::registry().with(filter);

    let _ = if logging.log_format.eq_ignore_ascii_case("json") {
        registry.with(fmt_layer.json()).try_init()
    } else {
        registry.with(fmt_layer).try_init()
    };

    if otel.otel_enable_telemetry {
        info!(
            endpoint = %otel.otel_exporter_otlp_endpoint,
            protocol = %otel.otel_exporter_otlp_protocol,
            sample_rate = otel.otel_tracing_sample_rate,
            "OTel telemetry enabled (exporter wiring deferred to D-AUTH phase)",
        );
    }
}

/// Shut down `OTel` pipeline on process exit. No-op in the scaffold;
/// real implementation will call `global::shutdown_tracer_provider()`.
pub fn shutdown_tracing() {}
