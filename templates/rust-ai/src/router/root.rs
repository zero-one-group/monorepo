//! `/` and `/health-check` routes.
//!
//! Mirrors `app/router/root.py` byte-for-byte on the response shapes:
//!
//! - `GET /` returns `200 {"message": "Welcome to the Machine Learning API"}`
//! - `GET /health-check` returns `200 {"status": "ok", "database": "connected"}`
//!   on success, or `503 {"detail": "Database connection error"}` on DB failure

use axum::Json;
use axum::extract::State;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use serde_json::json;

use crate::AppState;

pub fn router() -> axum::Router<AppState> {
    axum::Router::new()
        .route("/", get(root))
        .route("/health-check", get(health_check))
}

async fn root() -> Json<serde_json::Value> {
    Json(json!({"message": "Welcome to the Machine Learning API"}))
}

async fn health_check(State(state): State<AppState>) -> Response {
    tracing::info!("Performing health check...");
    if state.db.check_connection().await {
        tracing::info!("Health check successful: Database connection verified.");
        (
            StatusCode::OK,
            Json(json!({"status": "ok", "database": "connected"})),
        )
            .into_response()
    } else {
        tracing::error!("Health check failed: Database connection error.");
        (
            StatusCode::SERVICE_UNAVAILABLE,
            Json(json!({"detail": "Database connection error"})),
        )
            .into_response()
    }
}
