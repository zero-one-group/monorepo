//! `JwtGenerator` — port of `pkg/apputils/jwt_utils.go`.
//!
//! HS256-only (design 3.4: `JWT_ALGORITHM` config deleted).
//! Access token leeway: 30 seconds (fix audit §9.15).
//! Refresh token `aud` is hardcoded to `"client-app"` (design 3.6:
//! `X-App-Audience` plumbing deleted).
//!
//! **Claim shape is preserved byte-for-byte from the Go source**
//! so existing access tokens parse equivalently:
//!
//! Access: `{iss, iat, exp, sub, typ:"access", user_id, email, sid}`
//!         (no `aud`, no `jti`)
//! Refresh: `{iss, iat, exp, sub, aud:"client-app", typ:"refresh", jti}`

use std::time::Duration;

use chrono::{DateTime, Utc};
use jsonwebtoken::{Algorithm, DecodingKey, EncodingKey, Header, Validation, decode, encode};
use serde::{Deserialize, Serialize};

use crate::domain::AppError;

/// Hardcoded audience for refresh tokens (D-OPEN-2, design 3.6).
pub const REFRESH_TOKEN_AUDIENCE: &str = "client-app";

/// Access token payload — the app-specific claims flattened into
/// the JWT alongside the standard iss/iat/exp/sub/typ fields.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AccessTokenPayload {
    pub user_id: String,
    pub email: String,
    pub sid: String,
}

/// Decoded access token claims (what `verify_access` returns).
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AccessTokenClaims {
    pub iss: String,
    pub iat: i64,
    pub exp: i64,
    pub sub: String,
    pub typ: String,
    pub user_id: String,
    pub email: String,
    pub sid: String,
}

/// Decoded refresh token claims (what `verify_refresh` returns).
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RefreshTokenClaims {
    pub iss: String,
    pub iat: i64,
    pub exp: i64,
    pub sub: String,
    pub aud: String,
    pub typ: String,
    pub jti: String,
}

/// JWT generator + validator.
#[derive(Debug, Clone)]
pub struct JwtGenerator {
    secret: Vec<u8>,
    access_expiry: Duration,
    refresh_expiry: Duration,
    issuer: String,
}

impl JwtGenerator {
    pub fn new(
        secret: impl Into<Vec<u8>>,
        access_expiry: Duration,
        refresh_expiry: Duration,
        issuer: impl Into<String>,
    ) -> Self {
        Self {
            secret: secret.into(),
            access_expiry,
            refresh_expiry,
            issuer: issuer.into(),
        }
    }

    #[must_use]
    pub fn access_expiry(&self) -> Duration {
        self.access_expiry
    }

    #[must_use]
    pub fn refresh_expiry(&self) -> Duration {
        self.refresh_expiry
    }

    /// Sign an access token. Subject is the user id string;
    /// payload fields are flattened into the claims.
    pub fn sign_access(
        &self,
        subject: &str,
        payload: &AccessTokenPayload,
    ) -> Result<String, AppError> {
        let now = Utc::now();
        let exp = now + chrono::Duration::from_std(self.access_expiry).unwrap_or_default();

        let claims = AccessTokenClaims {
            iss: self.issuer.clone(),
            iat: now.timestamp(),
            exp: exp.timestamp(),
            sub: subject.to_string(),
            typ: "access".to_string(),
            user_id: payload.user_id.clone(),
            email: payload.email.clone(),
            sid: payload.sid.clone(),
        };

        encode(
            &Header::new(Algorithm::HS256),
            &claims,
            &EncodingKey::from_secret(&self.secret),
        )
        .map_err(|e| AppError::Internal(anyhow::anyhow!("sign access token: {e}")))
    }

    /// Sign a refresh token. `jti` should be the `UUIDv7` id of the
    /// `refresh_tokens` row (matches Go's pattern).
    pub fn sign_refresh(&self, subject: &str, jti: &str) -> Result<String, AppError> {
        let now = Utc::now();
        let exp = now + chrono::Duration::from_std(self.refresh_expiry).unwrap_or_default();

        let claims = RefreshTokenClaims {
            iss: self.issuer.clone(),
            iat: now.timestamp(),
            exp: exp.timestamp(),
            sub: subject.to_string(),
            aud: REFRESH_TOKEN_AUDIENCE.to_string(),
            typ: "refresh".to_string(),
            jti: jti.to_string(),
        };

        let mut header = Header::new(Algorithm::HS256);
        header.typ = Some("JWT".to_string());

        encode(&header, &claims, &EncodingKey::from_secret(&self.secret))
            .map_err(|e| AppError::Internal(anyhow::anyhow!("sign refresh token: {e}")))
    }

    /// Parse + validate an access token. Enforces signature, `exp`,
    /// and a 30-second clock skew leeway (fix audit §9.15).
    pub fn verify_access(&self, token: &str) -> Result<AccessTokenClaims, AppError> {
        let mut validation = Validation::new(Algorithm::HS256);
        validation.leeway = 30;
        // Access tokens have no `aud` — disable audience validation.
        validation.validate_aud = false;
        validation.set_required_spec_claims(&["exp", "iat", "sub"]);

        let data = decode::<AccessTokenClaims>(
            token,
            &DecodingKey::from_secret(&self.secret),
            &validation,
        )
        .map_err(|_| AppError::Unauthorized)?;

        if data.claims.typ != "access" {
            return Err(AppError::Unauthorized);
        }
        Ok(data.claims)
    }

    /// Parse + validate a refresh token. Enforces signature, `exp`,
    /// `typ == "refresh"`, and a 30-second leeway. Audience is
    /// extracted but NOT enforced (matches Go: audit §9.19 —
    /// latent design, not exploitable).
    pub fn verify_refresh(&self, token: &str) -> Result<RefreshTokenClaims, AppError> {
        let mut validation = Validation::new(Algorithm::HS256);
        validation.leeway = 30;
        validation.validate_aud = false;
        validation.set_required_spec_claims(&["exp", "iat", "sub"]);

        let data = decode::<RefreshTokenClaims>(
            token,
            &DecodingKey::from_secret(&self.secret),
            &validation,
        )
        .map_err(|_| AppError::Unauthorized)?;

        if data.claims.typ != "refresh" {
            return Err(AppError::Unauthorized);
        }
        Ok(data.claims)
    }

    /// Exposes the issuer string for external use (e.g., logging).
    #[must_use]
    pub fn issuer(&self) -> &str {
        &self.issuer
    }

    /// Helper: convert a `chrono::DateTime<Utc>` to a unix second.
    #[must_use]
    pub fn to_unix(ts: DateTime<Utc>) -> i64 {
        ts.timestamp()
    }
}
