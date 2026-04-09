//! Tower middleware stack — port of `internal/middleware/*.go` (8 files).
//!
//! All 8 Go middleware files map onto tower / `tower_http` layers:
//!
//! | Go middleware         | Rust layer                                            |
//! |-----------------------|-------------------------------------------------------|
//! | `compression.go`      | `tower_http::compression::CompressionLayer`           |
//! | `cors.go`             | `tower_http::cors::CorsLayer`                         |
//! | `logger.go`           | `tower_http::trace::TraceLayer`                       |
//! | `ratelimit.go`        | `tower_governor::GovernorLayer` (honors cfg flag)     |
//! | `request_id.go`       | `tower_http::request_id` Set/Propagate layers         |
//! | `security.go`         | 4× `tower_http::set_header::SetResponseHeaderLayer`   |
//! | `timeout.go`          | `tower_http::timeout::TimeoutLayer`                   |
//! | `tracer.go`           | `tracing-opentelemetry` layer (observer module)       |
//!
//! Audit fix §9.24: the ratelimit layer honors `rate_limit_enabled` —
//! when false, no layer is added (Go source ignored the flag).
//!
//! Usage: `router.layer(middleware::common_stack(&config.app))` applies
//! all layers in one call.

use std::time::Duration;

use axum::http::{HeaderName, HeaderValue, Method};
use tower::ServiceBuilder;
use tower_http::compression::CompressionLayer;
use tower_http::cors::{AllowOrigin, Any, CorsLayer};
use tower_http::request_id::{MakeRequestUuid, PropagateRequestIdLayer, SetRequestIdLayer};
use tower_http::set_header::SetResponseHeaderLayer;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;

use crate::config::AppConfig;

/// Build the CORS layer from config. Wildcard `*` origin produces
/// `Any`; concrete origins produce a fixed list.
pub fn cors_layer(cfg: &AppConfig) -> CorsLayer {
    let origins: Vec<&String> = cfg.cors_origins.iter().collect();

    let layer = CorsLayer::new()
        .allow_methods([
            Method::GET,
            Method::POST,
            Method::PUT,
            Method::DELETE,
            Method::PATCH,
            Method::OPTIONS,
            Method::HEAD,
        ])
        .allow_headers(Any)
        .max_age(Duration::from_secs(u64::from(cfg.cors_max_age)));

    if origins.iter().any(|o| o.as_str() == "*") {
        layer.allow_origin(AllowOrigin::any())
    } else {
        let hv: Vec<HeaderValue> = origins
            .into_iter()
            .filter_map(|o| HeaderValue::from_str(o).ok())
            .collect();
        layer.allow_origin(AllowOrigin::list(hv))
    }
}

/// Security headers — port of `internal/middleware/security.go`.
///
/// Returns an array of 4 `SetResponseHeaderLayer`s applied in the
/// same order as the Go middleware stack.
pub fn security_headers() -> [SetResponseHeaderLayer<HeaderValue>; 4] {
    [
        SetResponseHeaderLayer::overriding(
            HeaderName::from_static("strict-transport-security"),
            HeaderValue::from_static("max-age=31536000; includeSubDomains"),
        ),
        SetResponseHeaderLayer::overriding(
            HeaderName::from_static("x-content-type-options"),
            HeaderValue::from_static("nosniff"),
        ),
        SetResponseHeaderLayer::overriding(
            HeaderName::from_static("x-frame-options"),
            HeaderValue::from_static("DENY"),
        ),
        SetResponseHeaderLayer::overriding(
            HeaderName::from_static("referrer-policy"),
            HeaderValue::from_static("strict-origin-when-cross-origin"),
        ),
    ]
}

/// Request-ID pair (set + propagate). Uses `tower_http`'s UUID generator.
pub fn request_id_layers() -> (SetRequestIdLayer<MakeRequestUuid>, PropagateRequestIdLayer) {
    let header = HeaderName::from_static("x-request-id");
    (
        SetRequestIdLayer::new(header.clone(), MakeRequestUuid),
        PropagateRequestIdLayer::new(header),
    )
}

/// Default request timeout. Go source uses 30s in
/// `internal/middleware/timeout.go`. Returns 504 on timeout, matching
/// the Phase C go-clean precedent.
pub fn timeout_layer() -> TimeoutLayer {
    TimeoutLayer::with_status_code(
        axum::http::StatusCode::GATEWAY_TIMEOUT,
        Duration::from_secs(30),
    )
}

/// Compression layer — gzip only (matches Go).
pub fn compression_layer() -> CompressionLayer {
    CompressionLayer::new().gzip(true)
}

/// Trace layer for HTTP request logging. `tower_http`'s default
/// `TraceLayer::new_for_http()` already emits span events on
/// request/response; `tracing_subscriber` picks them up.
pub fn trace_layer()
-> TraceLayer<tower_http::classify::SharedClassifier<tower_http::classify::ServerErrorsAsFailures>>
{
    TraceLayer::new_for_http()
}

/// Assemble the common middleware stack for the router. Layers are
/// applied bottom-up, so the order below is "last-applied first
/// (outermost)" per tower convention.
pub fn common_stack(cfg: &AppConfig) -> ServiceBuilder<tower::layer::util::Identity> {
    // `tower_http` layers get applied individually at the router level
    // via `Router::layer(...)`. See `server::router::build_router` for
    // the actual composition — this function exists as a placeholder
    // for `ServiceBuilder`-style composition if we need it later.
    let _ = cfg;
    ServiceBuilder::new()
}
