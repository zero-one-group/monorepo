//! Auth handler: `POST /api/v1/auth/login`.
//!
//! Mirrors `apps/{{ package_name | kebab_case }}/internal/rest/auth.go`. Returns:
//! - `200 ResponseSingleData<LoginResponse>` on success
//! - `400 ResponseSingleData<Empty>` on payload bind failure
//! - `401 ResponseSingleData<Empty>` on any auth failure

use axum::Json;
use axum::extract::State;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};

use crate::AppState;
use crate::domain::auth::{LoginRequest, LoginResponse};
use crate::domain::response::{Empty, ResponseSingleData};

/// Custom extraction: Axum 0.8's `Json` extractor returns a 422 on
/// deserialization failure by default, but the Go handler returns 400
/// with a specific body shape. We handle that here by pulling bytes
/// then parsing manually.
pub async fn login(State(state): State<AppState>, body: axum::body::Bytes) -> Response {
    let Ok(req) = serde_json::from_slice::<LoginRequest>(&body) else {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "Invalid request payload",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    };

    if let Ok(login) = state.auth_service.login(&req.email, &req.password).await {
        let body = ResponseSingleData::<LoginResponse>::new(
            StatusCode::OK.as_u16(),
            login,
            "Successfully logged in",
        );
        (StatusCode::OK, Json(body)).into_response()
    } else {
        // The Go handler collapses every auth failure into a single
        // 401 with "Invalid email or password" — preserve exactly.
        let body = ResponseSingleData::new(
            StatusCode::UNAUTHORIZED.as_u16(),
            Empty::default(),
            "Invalid email or password",
        );
        (StatusCode::UNAUTHORIZED, Json(body)).into_response()
    }
}
