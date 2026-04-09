//! JWT validation middleware.
//!
//! Mirrors `apps/go-clean/internal/rest/middleware/jwt.go::ValidateUserToken`.
//! Reads the `Authorization` header, expects `Bearer <token>` format,
//! validates the token with `utils::jwt::validate_token`, stores the
//! decoded claim in request extensions, and calls `next.run(...)`.
//!
//! Error responses match Go byte-for-byte:
//! - Missing Authorization → 401 Unauthorized (Echo's default body)
//! - Malformed Authorization (not `Bearer X`) → 400 Bad Request
//!   `{"code":400,"data":{},"message":"invalid bearer token"}`
//! - Invalid token → 401 Unauthorized (Echo's default body)
//!
//! Echo's default 401 response is `{"message":"Unauthorized"}` (no
//! envelope). We replicate that by returning a plain JSON body on 401.

use axum::Json;
use axum::extract::{Request, State};
use axum::http::StatusCode;
use axum::middleware::Next;
use axum::response::{IntoResponse, Response};
use serde_json::json;

use crate::AppState;
use crate::domain::response::{Empty, ResponseSingleData};
use crate::utils::jwt::validate_token;

pub async fn require_auth(State(state): State<AppState>, mut req: Request, next: Next) -> Response {
    // Extract Authorization header
    let auth_header = req
        .headers()
        .get(http::header::AUTHORIZATION)
        .and_then(|h| h.to_str().ok());
    let Some(auth_header) = auth_header else {
        return echo_unauthorized();
    };

    // Split on whitespace, expect ["Bearer", "<token>"]
    let mut parts = auth_header.splitn(2, ' ');
    let scheme = parts.next().unwrap_or("");
    let token = parts.next().unwrap_or("");
    if scheme != "Bearer" || token.is_empty() {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "invalid bearer token",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    }

    // Validate with the same secret the auth service uses
    let Ok(claim) = validate_token(token, &state.env) else {
        return echo_unauthorized();
    };

    // Stash the claim in request extensions so handlers can read
    // `Extension<JwtClaim>` if they need the authenticated user.
    req.extensions_mut().insert(claim);

    next.run(req).await
}

fn echo_unauthorized() -> Response {
    // Echo's default Unauthorized response is `{"message":"Unauthorized"}`.
    (
        StatusCode::UNAUTHORIZED,
        Json(json!({"message":"Unauthorized"})),
    )
        .into_response()
}
