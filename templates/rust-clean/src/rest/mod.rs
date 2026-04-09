//! REST handlers + middleware + router composition.
//! Mirrors `apps/{{ package_name | kebab_case }}/internal/rest/`.

pub mod auth;
pub mod middleware;
pub mod root;
pub mod user;

use axum::Router;
use axum::middleware::from_fn_with_state;
use axum::routing::{get, post};
use http::{HeaderName, HeaderValue};
use tower_http::compression::CompressionLayer;
use tower_http::cors::{Any, CorsLayer};
use tower_http::request_id::{MakeRequestUuid, PropagateRequestIdLayer, SetRequestIdLayer};
use tower_http::set_header::SetResponseHeaderLayer;
use tower_http::timeout::TimeoutLayer;
use tower_http::trace::TraceLayer;

use crate::AppState;
use crate::rest::middleware::jwt::require_auth;

/// Build the full router with all middleware + route groups.
///
/// Middleware ordering mirrors `apps/{{ package_name | kebab_case }}/main.go` top-to-bottom:
/// 1. `RequestID` (injects `X-Request-ID`)
/// 2. Structured logger (`TraceLayer` → tracing spans per request)
/// 3. CORS (permissive)
/// 4. Security headers
/// 5. Compression (gzip)
/// 6. Rate limit (10 req/s, burst 20, per-IP) — wired in
///    `lib.rs::serve` via `into_make_service_with_connect_info` because
///    `tower_governor` needs the connect-info extractor for peer IPs.
/// 7. Timeout (30s)
pub fn build(state: AppState) -> Router {
    let request_id_header: HeaderName = HeaderName::from_static("x-request-id");

    let api_v1_users = Router::new()
        .route("/", get(user::list_users).post(user::create_user))
        .route(
            "/{id}",
            get(user::get_user)
                .put(user::update_user)
                .delete(user::delete_user),
        )
        .route_layer(from_fn_with_state(state.clone(), require_auth));

    let api_v1_auth = Router::new().route("/login", post(auth::login));

    let api_v1 = Router::new()
        .nest("/users", api_v1_users)
        .nest("/auth", api_v1_auth);

    let base = Router::new()
        .route("/", get(root::health))
        .nest("/api/v1", api_v1)
        .with_state(state)
        // Middleware (applied innermost → outermost in code; axum reverses
        // the order at runtime so the layer declared last runs first).
        .layer(TimeoutLayer::with_status_code(
            http::StatusCode::REQUEST_TIMEOUT,
            std::time::Duration::from_secs(30),
        ))
        .layer(CompressionLayer::new());

    apply_security_headers(base)
        .layer(
            CorsLayer::new()
                .allow_origin(Any)
                .allow_methods(Any)
                .allow_headers(Any),
        )
        .layer(TraceLayer::new_for_http())
        .layer(PropagateRequestIdLayer::new(request_id_header.clone()))
        .layer(SetRequestIdLayer::new(request_id_header, MakeRequestUuid))
}

/// Apply the 4 security headers that the Go middleware adds.
/// Factored as a function that wraps the router rather than a
/// pre-composed layer stack so the generic types stay readable.
fn apply_security_headers(router: Router) -> Router {
    router
        .layer(SetResponseHeaderLayer::if_not_present(
            HeaderName::from_static("x-content-type-options"),
            HeaderValue::from_static("nosniff"),
        ))
        .layer(SetResponseHeaderLayer::if_not_present(
            HeaderName::from_static("x-frame-options"),
            HeaderValue::from_static("DENY"),
        ))
        .layer(SetResponseHeaderLayer::if_not_present(
            HeaderName::from_static("referrer-policy"),
            HeaderValue::from_static("strict-origin-when-cross-origin"),
        ))
        .layer(SetResponseHeaderLayer::if_not_present(
            HeaderName::from_static("strict-transport-security"),
            HeaderValue::from_static("max-age=31536000; includeSubDomains"),
        ))
}
