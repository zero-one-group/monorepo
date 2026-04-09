//! `apputils/` — port of `pkg/apputils/` from the Go source.
//!
//! Submodules match the Go file names one-to-one:
//! - `jwt`        — `JwtGenerator` with HS256 + rotation helpers
//! - `password`   — argon2id PHC hasher with 16-byte salt (D-OPEN-2)
//! - `generator`  — URL-safe token generator (38 alpha + 10-digit ts)
//! - `validation` — validator error-map formatter
//!
//! Deliberately NOT ported from Go:
//! - `user_agent.rs` — unused by the corrected-port auth flow; IP /
//!   user-agent capture happens at the handler layer via axum's
//!   `ConnectInfo<SocketAddr>` and request headers directly.
//! - `HeadersContextKey` / `JWTClaimsContextKey` — deleted per
//!   design 3.6 (`X-App-Audience` plumbing removed).

pub mod generator;
pub mod jwt;
pub mod password;
pub mod validation;

pub use generator::{generate_url_safe_token, sha256_bytes, sha256_hex};
pub use jwt::{AccessTokenClaims, AccessTokenPayload, JwtGenerator, RefreshTokenClaims};
pub use password::PasswordHasher;
pub use validation::validation_errors_to_map;
