//! REST middleware. Mirrors `apps/{{ package_name | kebab_case }}/internal/rest/middleware/`.
//!
//! The only middleware that needs handler-level state is the JWT
//! validator — everything else (CORS, compression, timeout, request ID,
//! rate limit, security headers) is composed from generic tower layers
//! in `rest::build`.

pub mod jwt;
