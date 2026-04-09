//! Auth DB models — port of `modules/auth/models/model.go`.
//!
//! **D-OPEN-1 deltas vs Go source**:
//! - `Session::token_hash` is `Vec<u8>` (raw 32-byte SHA-256) instead
//!   of Go's `string` (hex-ASCII). Wire-shape JSON drift is accepted
//!   because session CRUD endpoints are never invoked externally by
//!   go-modular itself — only internal signin code creates sessions.
//! - `RefreshToken::token_hash` stays `Vec<u8>` but now holds the
//!   raw 32-byte digest instead of the hex-ASCII it did in Go.
//!
//! The `serde_bytes::ByteBuf`-style hex-string serialization is
//! NOT used. Rust sends raw byte arrays as JSON numeric arrays
//! (serde default). That's a minor drift from Go's base64 bytea
//! encoding but harmless — no client reads these fields directly.

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::types::Json;
use utoipa::ToSchema;
use uuid::Uuid;

use crate::modules::user::models::{User, UserMetadata as _UserMetadata};

// Re-export so consumers don't have to dance through the user module.
pub use crate::modules::user::models::UserMetadata;

// ----- user_passwords -----

/// Row shape for `public.user_passwords`. The `password_hash` column
/// is `BYTEA` and stores a PHC-encoded argon2id string (the bytes of
/// the ASCII string, not a raw key). Matches the Go source.
#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow, ToSchema)]
pub struct UserPassword {
    pub user_id: Uuid,
    #[serde(with = "phc_bytes")]
    #[schema(value_type = String, format = "byte")]
    pub password_hash: Vec<u8>,
    pub created_at: DateTime<Utc>,
    pub updated_at: Option<DateTime<Utc>>,
}

// Serialize the PHC bytes as a plain string (the PHC ASCII form).
mod phc_bytes {
    use serde::{Deserialize, Deserializer, Serializer};

    pub fn serialize<S: Serializer>(bytes: &[u8], s: S) -> Result<S::Ok, S::Error> {
        let phc = std::str::from_utf8(bytes).map_err(serde::ser::Error::custom)?;
        s.serialize_str(phc)
    }

    pub fn deserialize<'de, D: Deserializer<'de>>(d: D) -> Result<Vec<u8>, D::Error> {
        let s = String::deserialize(d)?;
        Ok(s.into_bytes())
    }
}

// ----- sessions -----

/// Row shape for `public.sessions`. Note D-OPEN-1 drift: `token_hash`
/// holds the raw 32-byte SHA-256 digest of the refresh JWT, not a
/// 64-char hex string like Go.
///
/// `ip_address` is stored as `Option<String>` in the Rust struct but
/// the DB column is `INET`. sqlx's default feature set does not
/// include the `ipnetwork` type, so we use explicit `$N::INET` casts
/// on INSERT and `ip_address::TEXT as ip_address` on SELECT. The
/// handler layer parses raw `String` from `X-Real-IP` / `X-Forwarded-For`
/// headers and binds it directly.
#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow, ToSchema)]
pub struct Session {
    pub id: Uuid,
    pub user_id: Uuid,
    #[schema(value_type = String, format = "byte")]
    pub token_hash: Vec<u8>,
    pub user_agent: Option<String>,
    pub device_name: Option<String>,
    pub device_fingerprint: Option<String>,
    pub ip_address: Option<String>,
    pub expires_at: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
    pub refreshed_at: Option<DateTime<Utc>>,
    pub revoked_at: Option<DateTime<Utc>>,
    pub revoked_by: Option<Uuid>,
}

// ----- refresh_tokens -----

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow, ToSchema)]
pub struct RefreshToken {
    pub id: Uuid,
    pub user_id: Uuid,
    pub session_id: Option<Uuid>,
    #[schema(value_type = String, format = "byte")]
    pub token_hash: Vec<u8>,
    pub ip_address: Option<String>,
    pub user_agent: Option<String>,
    pub expires_at: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
    pub revoked_at: Option<DateTime<Utc>>,
    pub revoked_by: Option<Uuid>,
}

// ----- one_time_tokens -----

/// `OneTimeTokenSubject` — enum for the `subject` column. Matches
/// the two Go constants plus reserves room for future additions.
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "snake_case")]
pub enum OneTimeTokenSubject {
    EmailOtp,
    EmailVerification,
}

impl OneTimeTokenSubject {
    #[must_use]
    pub fn as_str(self) -> &'static str {
        match self {
            Self::EmailOtp => "email_otp",
            Self::EmailVerification => "email_verification",
        }
    }
}

impl std::fmt::Display for OneTimeTokenSubject {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str((*self).as_str())
    }
}

/// Row shape for `public.one_time_tokens`. `token_hash` stays TEXT
/// (lowercase hex SHA-256) — not in D-OPEN-1's BYTEA harmonization
/// scope. `metadata` is JSONB.
#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow, ToSchema)]
pub struct OneTimeToken {
    pub id: Uuid,
    pub user_id: Option<Uuid>,
    pub subject: String,
    pub token_hash: String,
    pub relates_to: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    #[schema(value_type = Option<Object>)]
    pub metadata: Option<Json<serde_json::Value>>,
    pub created_at: DateTime<Utc>,
    pub expires_at: DateTime<Utc>,
    pub last_sent_at: Option<DateTime<Utc>>,
}

// ----- Response shapes -----

/// Response shape returned by signin and token-refresh endpoints.
/// Matches the Go `AuthenticatedUser` struct byte-for-byte.
#[derive(Debug, Clone, Serialize, ToSchema)]
pub struct AuthenticatedUser {
    pub user: User,
    pub access_token: String,
    pub refresh_token: String,
    pub session_id: Option<Uuid>,
    pub token_expiry: DateTime<Utc>,
}

/// `SignInResponse` — wrapper around `AuthenticatedUser` matching
/// the Go source. The Go struct embeds `AuthenticatedUser`, so the
/// JSON body is the same shape.
pub type SignInResponse = AuthenticatedUser;

// Silence the unused-import warning until the user module's
// UserMetadata is referenced by a repository call.
#[allow(dead_code)]
const _UNUSED_USER_METADATA: Option<_UserMetadata> = None;
