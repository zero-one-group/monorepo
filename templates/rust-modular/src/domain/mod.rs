//! Domain types shared across modules.
//!
//! - `error` — `AppError` enum with `IntoResponse` impl that matches
//!   {{ package_name | kebab_case }}'s `{"error": "..."}` shape.
//! - `response` — `MessageResponse` for `{"message": "..."}` payloads.

pub mod error;
pub mod response;

pub use error::AppError;
pub use response::MessageResponse;
