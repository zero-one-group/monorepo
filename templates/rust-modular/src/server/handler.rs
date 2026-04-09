//! Infrastructure route handlers: healthz, api-docs, openapi.json, 404.
//!
//! Matches the Go `internal/server/handler.go` route set. The api-docs
//! and openapi.json responses are placeholders until D-DOC-1 wires
//! utoipa's surface.

use axum::Json;
use axum::extract::State;
use axum::http::StatusCode;
use axum::response::IntoResponse;
use serde_json::json;
use utoipa::OpenApi;

use crate::AppState;

/// `GET /healthz` — liveness probe. Matches the `alexliesenfeld/health`
/// shape used by Go: `{"status": "up"}` on success with a per-check
/// breakdown. The only check is a database ping.
pub async fn healthz(State(state): State<AppState>) -> impl IntoResponse {
    let db_status = if sqlx::query("SELECT 1").execute(&state.pool).await.is_ok() {
        "up"
    } else {
        "down"
    };

    let overall = if db_status == "up" { "up" } else { "down" };
    let status_code = if overall == "up" {
        StatusCode::OK
    } else {
        StatusCode::SERVICE_UNAVAILABLE
    };

    (
        status_code,
        Json(json!({
            "status": overall,
            "details": {
                "database": { "status": db_status }
            }
        })),
    )
}

/// `GET /api-docs` — Scalar / Swagger UI. Placeholder for D-DOC-1.
pub async fn api_docs(State(state): State<AppState>) -> impl IntoResponse {
    if !state.config.app.enable_api_docs {
        return (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "API docs are disabled"})),
        );
    }
    (
        StatusCode::OK,
        Json(json!({
            "message": "API docs placeholder — full utoipa surface lands in D-DOC-1",
            "openapi_url": "/api/openapi.json"
        })),
    )
}

/// `GET /api/openapi.json` — generated `OpenAPI` 3.0 spec.
///
/// Built at compile time by utoipa from the per-handler
/// `#[utoipa::path]` annotations collected into
/// [`crate::openapi::ApiDoc`].
pub async fn openapi_json(State(state): State<AppState>) -> impl IntoResponse {
    if !state.config.app.enable_api_docs {
        return (
            StatusCode::NOT_FOUND,
            Json(json!({"error": "API docs are disabled"})),
        );
    }

    let doc = crate::openapi::ApiDoc::openapi();
    match serde_json::to_value(&doc) {
        Ok(value) => (StatusCode::OK, Json(value)),
        Err(err) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(json!({"error": format!("serialize openapi: {err}")})),
        ),
    }
}

/// Catch-all 404.
pub async fn not_found() -> impl IntoResponse {
    (
        StatusCode::NOT_FOUND,
        Json(json!({"error": "route not found"})),
    )
}
