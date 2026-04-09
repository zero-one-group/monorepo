//! HTTP router composition.
//!
//! Mirrors `app/router/` — the root router exposes `GET /` and
//! `GET /health-check`, the openai router exposes `GET /openai/greetings`.
//! Both are merged into a single axum `Router` together with the CORS
//! and request-id middleware that the Python service used.

pub mod openai;
pub mod root;

use axum::Router;
use http::HeaderName;
use tower_http::cors::{Any, CorsLayer};
use tower_http::request_id::{MakeRequestUuid, PropagateRequestIdLayer, SetRequestIdLayer};

use crate::AppState;

/// Construct the application router with all middleware applied.
///
/// Middleware ordering matches `app/main.py`: a permissive CORS layer
/// (allow any origin / methods / headers) wrapped around a request-id
/// layer that injects a uuid v4 per inbound request.
pub fn build(state: AppState) -> Router {
    let request_id_header: HeaderName = HeaderName::from_static("x-request-id");

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    Router::new()
        .merge(root::router())
        .merge(openai::router())
        .with_state(state)
        .layer(PropagateRequestIdLayer::new(request_id_header.clone()))
        .layer(SetRequestIdLayer::new(request_id_header, MakeRequestUuid))
        .layer(cors)
}
