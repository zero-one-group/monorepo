//! Application error type.
//!
//! Mirrors the Go package-level errors in `domain/error.go`:
//! - `ErrInternalServerError`
//! - `ErrNotFound`
//! - `ErrConflict`
//! - `ErrBadParamInput`
//! - `ErrUserNotFound`
//!
//! Plus catches sqlx / bcrypt / JWT errors and converts them into HTTP
//! responses via `IntoResponse`. Every variant maps to a specific HTTP
//! status + a `ResponseSingleData<Empty>` envelope so the JSON shape
//! on the wire matches the Go service byte-for-byte.

use axum::Json;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use thiserror::Error;

use crate::domain::response::{Empty, ResponseSingleData};

#[derive(Debug, Error)]
pub enum AppError {
    #[error("internal Server Error")]
    InternalServerError,

    #[error("your requested Item is not found")]
    NotFound,

    #[error("your Item already exist")]
    Conflict,

    #[error("given Param is not valid")]
    BadParamInput,

    #[error("user not found")]
    UserNotFound,

    #[error("invalid email or password")]
    Unauthorized,

    #[error("invalid user ID format")]
    InvalidUserIdFormat,

    #[error("invalid request payload")]
    InvalidPayload,

    #[error("invalid bearer token")]
    InvalidBearer,

    #[error("database error: {0}")]
    Database(#[from] sqlx::Error),

    #[error("bcrypt error: {0}")]
    Bcrypt(#[from] bcrypt::BcryptError),

    #[error("jwt error: {0}")]
    Jwt(#[from] jsonwebtoken::errors::Error),

    #[error("{0}")]
    Other(String),
}

impl AppError {
    pub fn status_code(&self) -> StatusCode {
        match self {
            Self::NotFound | Self::UserNotFound => StatusCode::NOT_FOUND,
            Self::Conflict => StatusCode::CONFLICT,
            Self::BadParamInput | Self::InvalidUserIdFormat | Self::InvalidPayload => {
                StatusCode::BAD_REQUEST
            }
            Self::InvalidBearer => StatusCode::BAD_REQUEST,
            Self::Unauthorized => StatusCode::UNAUTHORIZED,
            Self::InternalServerError
            | Self::Database(_)
            | Self::Bcrypt(_)
            | Self::Jwt(_)
            | Self::Other(_) => StatusCode::INTERNAL_SERVER_ERROR,
        }
    }

    pub fn message(&self) -> String {
        self.to_string()
    }
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let status = self.status_code();
        let body = ResponseSingleData::new(status.as_u16(), Empty::default(), self.message());
        (status, Json(body)).into_response()
    }
}
