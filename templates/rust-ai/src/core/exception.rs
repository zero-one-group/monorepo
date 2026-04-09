//! Application error type.
//!
//! Mirrors `app/core/exception.py::AppError`. Deliberately matches the
//! Python error envelope so the JSON returned by the `AppError` handler is
//! identical:
//!
//! ```json
//! {
//!     "success": false,
//!     "message": "...",
//!     "error_code": "...",
//!     "data": null
//! }
//! ```

use http::StatusCode;
use serde_json::Value;
use thiserror::Error;

#[derive(Debug, Error)]
#[error("{message}")]
pub struct AppError {
    pub message: String,
    pub status_code: StatusCode,
    pub code: String,
    pub data: Option<Value>,
}

impl AppError {
    pub fn new(
        message: impl Into<String>,
        status_code: StatusCode,
        code: impl Into<String>,
    ) -> Self {
        Self {
            message: message.into(),
            status_code,
            code: code.into(),
            data: None,
        }
    }

    /// Default constructor — `400 BAD_REQUEST` with the supplied message.
    /// Mirrors Python's `AppError(message="Operation failed")`.
    pub fn bad_request(message: impl Into<String>) -> Self {
        Self::new(message, StatusCode::BAD_REQUEST, "BAD_REQUEST")
    }

    /// Mark this `AppError` with arbitrary additional data.
    #[must_use]
    pub fn with_data(mut self, data: Value) -> Self {
        self.data = Some(data);
        self
    }
}
