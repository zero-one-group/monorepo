//! Auth module route registration.
//!
//! Builds two routers: `public_routes()` (unauthenticated) and
//! `protected_routes()` (wrapped with `require_auth` middleware).
//!
//! The top-level `server/router.rs` nests both under `/api/v1/auth`.
//!
//! **Endpoint matrix** (14 total — 9 public, 5 protected):
//!
//! Public:
//! - POST /signin/email
//! - POST /signin/username
//! - GET  /verify-email
//! - POST /token/refresh   (NEW — replaces 4 deleted CRUD endpoints)
//! - POST /verification/email/initiate
//! - POST /verification/email/validate
//! - POST /password                (open in Go; Phase D enforces ownership)
//! - POST /session                 (keep-track CRUD)
//! - GET  /session/{sessionId}     (keep-track CRUD)
//!
//! Protected (JWT + session-check via `require_auth`):
//! - PUT    /password/{userId}
//! - PUT    /session
//! - DELETE /session/{sessionId}
//! - POST   /verification/email/revoke
//! - POST   /verification/email/resend

use axum::Router;
use axum::middleware::from_fn_with_state;
use axum::routing::{delete, get, post, put};

use crate::AppState;

use super::handler;
use super::middleware::require_auth;

/// Public (unauthenticated) routes.
pub fn public_routes() -> Router<AppState> {
    Router::new()
        .route("/signin/email", post(handler::sign_in_with_email))
        .route("/signin/username", post(handler::sign_in_with_username))
        .route("/verify-email", get(handler::verify_email_by_link))
        .route("/token/refresh", post(handler::rotate_refresh_token))
        .route(
            "/verification/email/initiate",
            post(handler::initiate_email_verification),
        )
        .route(
            "/verification/email/validate",
            post(handler::validate_email_verification),
        )
        .route("/password", post(handler::set_password))
        .route("/session", post(handler::create_session))
        .route("/session/{sessionId}", get(handler::get_session))
}

/// Protected routes (JWT-authenticated + session-check).
pub fn protected_routes(state: AppState) -> Router<AppState> {
    Router::new()
        .route("/password/{userId}", put(handler::update_password))
        .route("/session", put(handler::update_session))
        .route("/session/{sessionId}", delete(handler::delete_session))
        .route(
            "/verification/email/revoke",
            post(handler::revoke_email_verification),
        )
        .route(
            "/verification/email/resend",
            post(handler::resend_email_verification),
        )
        .route_layer(from_fn_with_state(state, require_auth))
}
