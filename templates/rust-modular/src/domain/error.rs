//! `AppError` — enum mapping every Phase D error condition to an
//! HTTP status + `{"error": "..."}` body shape that matches the Go
//! source's actual response convention.
//!
//! Variants specific to the Phase D corrected-port design:
//!
//! - `InvalidCredentials`        — 401 signin failure
//! - `EmailNotVerified`          — 401 signin blocked
//! - `RefreshTokenReuse`         — 401 + revoke-all (design 3.1)
//! - `SessionRevoked`            — 401 middleware check failure (3.2)
//! - `ConcurrentRefresh`         — 409 row-lock timeout (3.1)
//! - `OwnershipViolation`        — 403 set/update password authZ (3.9)
//! - `VerificationCooldown`      — 429 with Retry-After header (3.8)

use axum::Json;
use axum::http::{HeaderValue, StatusCode};
use axum::response::{IntoResponse, Response as AxumResponse};
use serde_json::{Value, json};
use thiserror::Error;

#[derive(Debug, Error)]
pub enum AppError {
    #[error("internal server error")]
    Internal(#[source] anyhow::Error),

    #[error(transparent)]
    Database(#[from] sqlx::Error),

    #[error("{0}")]
    NotFound(String),

    #[error("{0}")]
    BadRequest(String),

    #[error("{0}")]
    Conflict(String),

    #[error("Invalid email or password")]
    InvalidCredentials,

    #[error("Email is not verified")]
    EmailNotVerified,

    #[error("Unauthorized")]
    Unauthorized,

    #[error("Invalid bearer token")]
    InvalidBearer,

    #[error("Refresh token reuse detected — all sessions revoked")]
    RefreshTokenReuse,

    #[error("Session revoked")]
    SessionRevoked,

    #[error("Concurrent refresh in progress")]
    ConcurrentRefresh,

    #[error("Cannot modify another user's resource")]
    OwnershipViolation,

    #[error("Verification email cooldown: retry after {retry_after} seconds")]
    VerificationCooldown { retry_after: u64 },

    /// Validation failure with a per-field details map (matches the
    /// `apputils.ValidationErrorsToMap` shape in Go).
    #[error("Validation failed")]
    Validation(Value),
}

impl AppError {
    /// Map each variant to its HTTP status code.
    pub fn status(&self) -> StatusCode {
        match self {
            Self::Internal(_) | Self::Database(_) => StatusCode::INTERNAL_SERVER_ERROR,
            Self::NotFound(_) => StatusCode::NOT_FOUND,
            Self::BadRequest(_) | Self::InvalidBearer | Self::Validation(_) => {
                StatusCode::BAD_REQUEST
            }
            Self::Conflict(_) | Self::ConcurrentRefresh => StatusCode::CONFLICT,
            Self::InvalidCredentials
            | Self::EmailNotVerified
            | Self::Unauthorized
            | Self::RefreshTokenReuse
            | Self::SessionRevoked => StatusCode::UNAUTHORIZED,
            Self::OwnershipViolation => StatusCode::FORBIDDEN,
            Self::VerificationCooldown { .. } => StatusCode::TOO_MANY_REQUESTS,
        }
    }
}

impl From<anyhow::Error> for AppError {
    fn from(err: anyhow::Error) -> Self {
        Self::Internal(err)
    }
}

impl IntoResponse for AppError {
    fn into_response(self) -> AxumResponse {
        let status = self.status();
        let body = match &self {
            Self::Validation(details) => json!({
                "error": "Validation failed",
                "details": details,
            }),
            Self::VerificationCooldown { retry_after } => json!({
                "error": self.to_string(),
                "retry_after": retry_after,
            }),
            _ => json!({ "error": self.to_string() }),
        };

        let mut response = (status, Json(body)).into_response();

        // Retry-After header for verification cooldown (HTTP 429).
        if let Self::VerificationCooldown { retry_after } = self
            && let Ok(v) = HeaderValue::from_str(&retry_after.to_string())
        {
            response.headers_mut().insert("retry-after", v);
        }

        response
    }
}
