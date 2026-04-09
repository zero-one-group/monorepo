//! Auth request DTOs — port of `modules/auth/models/schema.go`.
//!
//! Phase D deletions vs Go source:
//! - `CreateRefreshTokenRequest` / `UpdateRefreshTokenRequest` are
//!   NOT ported. The 4 naive CRUD endpoints are removed per
//!   D-OPEN-5; only the new rotation endpoint (`TokenRefreshRequest`)
//!   replaces them.
//!
//! Validator attributes use `validator 0.20.0` derive syntax.
//! UUID + datetime fields are parsed in the handler layer with
//! `Uuid::parse_str` / `DateTime::parse_from_rfc3339` rather than
//! via validator custom functions — simpler and matches the
//! existing user module pattern.

use serde::{Deserialize, Serialize};
use utoipa::{IntoParams, ToSchema};
use validator::Validate;

// ----- Password -----

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct SetPasswordRequest {
    #[validate(length(min = 1))]
    pub user_id: String,

    #[validate(length(min = 8, message = "Minimum length is 8"))]
    pub password: String,

    #[validate(must_match(other = "password", message = "Must match password"))]
    pub password_confirmation: String,
}

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct UpdatePasswordRequest {
    #[validate(length(min = 8, message = "Minimum length is 8"))]
    pub current_password: String,

    #[validate(length(min = 8, message = "Minimum length is 8"))]
    pub new_password: String,

    #[validate(must_match(other = "new_password", message = "Must match new_password"))]
    pub password_confirmation: String,
}

// ----- Session CRUD -----

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct CreateSessionRequest {
    #[validate(length(min = 1))]
    pub user_id: String,

    #[validate(length(min = 1))]
    pub token_hash: String,

    pub user_agent: Option<String>,
    pub device_name: Option<String>,
    pub device_fingerprint: Option<String>,
    pub ip_address: Option<String>,

    #[validate(length(min = 1))]
    pub expires_at: String,
}

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct UpdateSessionRequest {
    #[validate(length(min = 1))]
    pub session_id: String,

    pub user_agent: Option<String>,
    pub device_name: Option<String>,
    pub device_fingerprint: Option<String>,
    pub ip_address: Option<String>,
    pub refreshed_at: Option<String>,
    pub revoked_at: Option<String>,
    pub revoked_by: Option<String>,
}

// ----- Signin -----

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct SignInWithEmailRequest {
    #[validate(email(message = "Must be a valid email address"))]
    pub email: String,

    #[validate(length(min = 1, message = "The password field is required"))]
    pub password: String,
}

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct SignInWithUsernameRequest {
    #[validate(length(min = 1, message = "The username field is required"))]
    pub username: String,

    #[validate(length(min = 1, message = "The password field is required"))]
    pub password: String,
}

// ----- NEW: Token rotation (replaces 4 deleted CRUD endpoints) -----

/// `POST /api/v1/auth/token/refresh` body.
///
/// This is the NEW endpoint that replaces the 4 naive refresh-token
/// CRUD endpoints deleted per D-OPEN-5. It implements
/// rotate-on-refresh with reuse detection per design 3.1.
#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct TokenRefreshRequest {
    #[validate(length(min = 1, message = "The refresh_token field is required"))]
    pub refresh_token: String,
}

// ----- Verification -----

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct InitiateEmailVerificationRequest {
    #[validate(email(message = "Must be a valid email address"))]
    pub email: String,

    #[serde(default)]
    #[validate(url(message = "Must be a valid URL"))]
    pub redirect_to: Option<String>,
}

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct ValidateEmailVerificationRequest {
    #[validate(length(min = 1, message = "The token field is required"))]
    pub token: String,
}

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct RevokeEmailVerificationRequest {
    #[validate(length(min = 1, message = "The token field is required"))]
    pub token: String,
}

#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct ResendEmailVerificationRequest {
    #[validate(email(message = "Must be a valid email address"))]
    pub email: String,
}

// ----- Query params -----

/// `GET /api/v1/auth/verify-email?token=...&redirect_to=...`.
#[derive(Debug, Clone, Deserialize, IntoParams)]
#[into_params(parameter_in = Query)]
pub struct VerifyEmailLinkQuery {
    pub token: String,
    #[serde(default)]
    pub redirect_to: Option<String>,
}

// ----- Response bodies -----

/// Simple success response for neutral endpoints (verification
/// initiate, resend).
#[derive(Debug, Clone, Serialize, ToSchema)]
pub struct NeutralResponse {
    pub message: String,
}
