//! URL-safe token generator — port of
//! `pkg/apputils/generator.go::GenerateURLSafeToken` plus SHA-256
//! hashing helpers used across the auth flow.
//!
//! Token shape matches Go exactly: `(length - 10)` random alphanumeric
//! characters followed by a 10-digit unix timestamp. The Go default
//! `length = 48` produces 38 random chars + 10 digits.
//!
//! **Randomness source:** we use `argon2::password_hash::rand_core::OsRng`
//! instead of pulling in the full `rand` crate. `OsRng` is the
//! operating-system CSPRNG, which is what the Go source uses via
//! `crypto/rand`. Avoiding a direct `rand` workspace dep keeps the
//! Phase D delta minimal (argon2 already transitively depends on
//! `rand_core`).

use std::fmt::Write as _;

use argon2::password_hash::rand_core::{OsRng, RngCore};
use sha2::{Digest, Sha256};

use crate::domain::AppError;

const ALNUM: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

/// Generate a URL-safe alphanumeric token of the given total length.
///
/// Last 10 characters are the current unix timestamp as decimal
/// digits. The remaining `(length - 10)` characters are ASCII
/// alphanumeric from `OsRng`.
///
/// Returns `AppError::Internal` if `length < 11` (no room for the
/// 10-digit timestamp plus at least one random char).
pub fn generate_url_safe_token(length: usize) -> Result<String, AppError> {
    if length < 11 {
        return Err(AppError::Internal(anyhow::anyhow!(
            "token length must be >= 11 (got {length})"
        )));
    }
    let token_len = length - 10;

    let mut out = String::with_capacity(length);
    let mut buf = [0u8; 64];
    while out.len() < token_len {
        OsRng.fill_bytes(&mut buf);
        for &b in &buf {
            if out.len() == token_len {
                break;
            }
            // Unbiased-ish: use the byte modulo 62. The bias is
            // (256 % 62) / 256 ~ 3%, matching Go's base64-then-filter
            // pattern which also isn't perfectly unbiased.
            let idx = (b as usize) % ALNUM.len();
            out.push(ALNUM[idx] as char);
        }
    }

    let unix_ts = chrono::Utc::now().timestamp();
    // Unix timestamps are 10 digits for the foreseeable future
    // (until ~2286), matching Go's `fmt.Sprintf("%s%d", ...)`.
    write!(&mut out, "{unix_ts}").map_err(|e| AppError::Internal(anyhow::anyhow!("{e}")))?;
    Ok(out)
}

/// SHA-256 of the given input, returned as a raw 32-byte array.
/// Used by the D-OPEN-1 harmonized `sessions.token_hash` and
/// `refresh_tokens.token_hash` columns.
#[must_use]
pub fn sha256_bytes(input: &[u8]) -> [u8; 32] {
    let mut hasher = Sha256::new();
    hasher.update(input);
    hasher.finalize().into()
}

/// SHA-256 of the given input, returned as a lowercase hex string.
/// Used by `one_time_tokens.token_hash` (which stays TEXT — not in
/// D-OPEN-1's BYTEA harmonization scope).
#[must_use]
pub fn sha256_hex(input: &[u8]) -> String {
    let bytes = sha256_bytes(input);
    let mut out = String::with_capacity(64);
    for b in bytes {
        let _ = write!(out, "{b:02x}");
    }
    out
}

#[cfg(test)]
mod tests {
    use std::fmt::Write as _;

    use super::{generate_url_safe_token, sha256_bytes, sha256_hex};

    #[test]
    fn token_has_requested_length_and_trailing_digits() {
        let token = generate_url_safe_token(48).unwrap();
        assert_eq!(token.len(), 48);
        // Last 10 chars must be digits (unix timestamp).
        let tail = &token[38..];
        assert!(tail.chars().all(|c| c.is_ascii_digit()), "tail: {tail}");
        // Remaining must be alphanumeric.
        let head = &token[..38];
        assert!(head.chars().all(|c| c.is_ascii_alphanumeric()));
    }

    #[test]
    fn token_rejects_too_short() {
        assert!(generate_url_safe_token(10).is_err());
    }

    #[test]
    fn sha256_hex_matches_bytes() {
        let bytes = sha256_bytes(b"hello");
        let hex = sha256_hex(b"hello");
        assert_eq!(hex.len(), 64);
        // Encode bytes to hex manually to verify consistency.
        let manual = bytes.iter().fold(String::with_capacity(64), |mut acc, b| {
            let _ = write!(acc, "{b:02x}");
            acc
        });
        assert_eq!(hex, manual);
    }
}
