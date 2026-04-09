//! `require_auth` — session-check JWT middleware (D-AUTH-14).
//!
//! Flow:
//! 1. Extract `Authorization: Bearer <token>` header.
//! 2. Verify JWT signature + exp with a 30-second leeway.
//! 3. Extract `sid` claim from the access token.
//! 4. `SELECT revoked_at, expires_at FROM sessions WHERE id = $1`.
//! 5. If row missing OR `revoked_at IS NOT NULL` OR `expires_at < NOW()`,
//!    return 401.
//! 6. Stash `caller_id` + `sid` in request extensions.
//!
//! **v1 has no cache** — a DB lookup happens on every protected
//! request. D-IT-4 will bench the p99 overhead; if it exceeds 5 ms,
//! a follow-up adds a `moka` LRU cache.

use axum::extract::{Request, State};
use axum::http::header::AUTHORIZATION;
use axum::middleware::Next;
use axum::response::Response;
use chrono::Utc;
use uuid::Uuid;

use crate::AppState;
use crate::domain::AppError;

/// Request-extension entry: `(caller_user_id, session_id)`.
#[derive(Debug, Clone, Copy)]
pub struct AuthContext {
    pub user_id: Uuid,
    pub session_id: Uuid,
}

pub async fn require_auth(
    State(state): State<AppState>,
    mut req: Request,
    next: Next,
) -> Result<Response, AppError> {
    // 1. Extract bearer header.
    let bearer = req
        .headers()
        .get(AUTHORIZATION)
        .and_then(|h| h.to_str().ok())
        .ok_or(AppError::Unauthorized)?;

    let token = bearer
        .strip_prefix("Bearer ")
        .ok_or(AppError::InvalidBearer)?
        .trim();
    if token.is_empty() {
        return Err(AppError::InvalidBearer);
    }

    // 2. Verify JWT signature + exp (with 30s leeway inside verify_access).
    let claims = state.auth_service.jwt().verify_access(token)?;

    // 3. Parse caller + session ids.
    let user_id = Uuid::parse_str(&claims.user_id).map_err(|_| AppError::Unauthorized)?;
    let session_id = Uuid::parse_str(&claims.sid).map_err(|_| AppError::Unauthorized)?;

    // 4. Session-row check (the fix to audit §9.3).
    let (revoked_at, expires_at) = state
        .auth_service
        .repository()
        .validate_session(session_id)
        .await?;
    if revoked_at.is_some() {
        return Err(AppError::SessionRevoked);
    }
    if expires_at < Utc::now() {
        return Err(AppError::SessionRevoked);
    }

    // 5. Stash context in request extensions for downstream handlers.
    req.extensions_mut().insert(AuthContext {
        user_id,
        session_id,
    });

    Ok(next.run(req).await)
}
