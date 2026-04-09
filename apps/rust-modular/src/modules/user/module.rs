//! User module route registration.
//!
//! Port of `modules/user/module.go` with the difference that the
//! JWT-protected middleware layer is NOT applied here — it lands
//! in D-AUTH-14 as `require_auth`, and the router in
//! `server::router::build_router` wraps the whole user surface in
//! `.route_layer(from_fn_with_state(state, require_auth))` once
//! that middleware exists.

use axum::Router;
use axum::routing::{delete, get, post, put};

use crate::AppState;

use super::handler;

/// Build the user module router. Paths are relative to the `/users`
/// prefix that the caller mounts under `/api/v1`.
pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/", post(handler::create_user))
        .route("/", get(handler::list_users))
        .route("/{userId}", get(handler::get_user))
        .route("/{userId}", put(handler::update_user))
        .route("/{userId}", delete(handler::delete_user))
}
