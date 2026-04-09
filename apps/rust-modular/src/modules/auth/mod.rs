//! Auth module — 14 endpoints (fix-track).
//!
//! Port of `apps/go-modular/modules/auth/` with the 8 corrected-port
//! fixes listed in the Phase D plan. Module layout:
//!
//! - `models` — DB model structs for the 4 auth tables plus
//!   `AuthenticatedUser` / `SignInResponse` response shapes.
//! - `schema` — Request DTOs for every endpoint.
//! - `repository` — 4 sqlx repositories (password, session,
//!   `refresh_token`, `one_time_token`). D-OPEN-1 BYTEA harmonization,
//!   drops the generic `update` on `refresh_tokens` per D-OPEN-5.
//! - `service` — 5 services (password, verification, session, signin,
//!   `token_refresh`). All 8 corrected-port fixes live here.
//! - `handler` — 14 HTTP handlers matching the endpoint matrix.
//! - `middleware` — `require_auth` session-check middleware.
//! - `module` — Route registration.

pub mod handler;
pub mod middleware;
pub mod models;
pub mod module;
pub mod repository;
pub mod schema;
pub mod service;

pub use middleware::require_auth;
pub use module::{protected_routes, public_routes};
pub use service::AuthService;
