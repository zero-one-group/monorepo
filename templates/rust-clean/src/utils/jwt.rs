//! JWT helper. Mirrors `apps/{{ package_name | kebab_case }}/utils/jwt.go`.
//!
//! - Access token: HS256-signed `JwtClaim` with `exp` = now + TTL minutes
//! - Refresh token: HS256-signed `RefreshClaim` with `exp` = now + 24h
//!
//! The refresh claim shape in the Go source is identical to the access
//! claim (both `{id, email, exp, iat}`), so we reuse a single struct.

use chrono::{Duration, Utc};
use jsonwebtoken::{Algorithm, DecodingKey, EncodingKey, Header, Validation, decode, encode};

use crate::config::env::Env;
use crate::domain::auth::JwtClaim;
use crate::domain::error::AppError;

/// Generate (`access_token`, `refresh_token`) using HS256 signing.
pub fn generate_token_pair(
    user_id: &str,
    email: &str,
    env: &Env,
) -> Result<(String, String), AppError> {
    let secret = env.jwt_secret.as_bytes();
    let key = EncodingKey::from_secret(secret);
    let header = Header::new(Algorithm::HS256);
    let now = Utc::now();

    // Access token (configured TTL, default 60 minutes).
    let access_claims = JwtClaim {
        id: user_id.to_owned(),
        email: email.to_owned(),
        iat: now.timestamp(),
        exp: (now + Duration::minutes(env.auth_token_expiry_minutes)).timestamp(),
    };
    let access_token = encode(&header, &access_claims, &key)?;

    // Refresh token (hardcoded 24h window, matching the Go source).
    let refresh_claims = JwtClaim {
        id: user_id.to_owned(),
        email: email.to_owned(),
        iat: now.timestamp(),
        exp: (now + Duration::hours(24)).timestamp(),
    };
    let refresh_token = encode(&header, &refresh_claims, &key)?;

    Ok((access_token, refresh_token))
}

/// Decode + validate an HS256 token. Returns the claim payload.
pub fn validate_token(token: &str, env: &Env) -> Result<JwtClaim, AppError> {
    let key = DecodingKey::from_secret(env.jwt_secret.as_bytes());
    let validation = Validation::new(Algorithm::HS256);
    let data = decode::<JwtClaim>(token, &key, &validation)?;
    Ok(data.claims)
}
