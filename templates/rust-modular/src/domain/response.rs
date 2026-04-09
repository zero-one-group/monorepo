//! Response envelope shapes used by {{ package_name | kebab_case }}.
//!
//! Unlike Phase C go-clean (which wrapped every response in a
//! `{code, data, message}` envelope), {{ package_name | kebab_case }} returns:
//! - **Success**: raw struct JSON (no wrapping envelope)
//! - **Updates / deletes**: `{"message": "..."}`
//! - **Errors**: `{"error": "..."}` plus optional `{"details": ...}`
//!
//! The user handler (`modules/user/handler.go`) makes this explicit:
//! ```text
//! return c.JSON(http.StatusCreated, user)  // raw struct
//! return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
//! return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
//! ```
//!
//! This helper struct gives us a type-safe call site for the
//! message-only shape used by update/delete success responses.

use serde::{Deserialize, Serialize};
use utoipa::ToSchema;

/// `{"message": "..."}` — used by update/delete success responses.
#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct MessageResponse {
    pub message: String,
}

/// Generic `{"error": "..."}` body returned by all error paths.
/// Not used in Rust code directly (handlers return `AppError` which
/// serializes to this shape via `IntoResponse`) but declared here
/// so it appears as a reusable schema in the `OpenAPI` spec.
#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct ErrorBody {
    pub error: String,
}

impl MessageResponse {
    pub fn new(message: impl Into<String>) -> Self {
        Self {
            message: message.into(),
        }
    }
}
