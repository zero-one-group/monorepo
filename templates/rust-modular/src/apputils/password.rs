//! Argon2id password hasher — port of `pkg/apputils/password.go`.
//!
//! **D-OPEN-2 delta vs Go source:**
//! - Salt length bumped from 8 bytes to 16 bytes (OWASP-recommended).
//! - All other argon2id params match Go exactly: 16 MiB memory,
//!   4 iterations, 2 parallelism, 32-byte key output.
//! - PHC string format is identical; old Go-era hashes (8-byte salt)
//!   still verify correctly because PHC encodes the salt length
//!   per-hash — the verifier reads params from the stored PHC, not
//!   from the hasher config.

use argon2::password_hash::rand_core::OsRng;
use argon2::password_hash::{PasswordHash, PasswordHasher as _, PasswordVerifier, SaltString};
use argon2::{Algorithm, Argon2, Params, Version};

use crate::domain::AppError;

/// Argon2id password hasher with fixed OWASP-compliant parameters.
#[derive(Debug, Clone)]
pub struct PasswordHasher {
    argon2: Argon2<'static>,
}

impl Default for PasswordHasher {
    fn default() -> Self {
        Self::new()
    }
}

impl PasswordHasher {
    /// Build with the Phase D default parameters (16 MiB / 4 iter /
    /// 2 parallel / 32-byte output).
    #[must_use]
    pub fn new() -> Self {
        let params = Params::new(16 * 1024, 4, 2, Some(32)).expect("valid argon2 params");
        let argon2 = Argon2::new(Algorithm::Argon2id, Version::V0x13, params);
        Self { argon2 }
    }

    /// Hash a password, returning a PHC string.
    ///
    /// Salt is 16 random bytes (D-OPEN-2). The `argon2` crate's
    /// `SaltString::generate` produces the OWASP-recommended length
    /// by default.
    pub fn hash(&self, password: &str) -> Result<String, AppError> {
        let salt = SaltString::generate(&mut OsRng);
        self.argon2
            .hash_password(password.as_bytes(), &salt)
            .map(|hash| hash.to_string())
            .map_err(|e| AppError::Internal(anyhow::anyhow!("argon2 hash: {e}")))
    }

    /// Verify a password against a stored PHC string. Returns
    /// `Ok(true)` on match, `Ok(false)` on mismatch, `Err` on a
    /// parse or crypto failure.
    pub fn verify(&self, password: &str, phc: &str) -> Result<bool, AppError> {
        let parsed = PasswordHash::new(phc)
            .map_err(|e| AppError::Internal(anyhow::anyhow!("parse argon2 phc: {e}")))?;
        match self.argon2.verify_password(password.as_bytes(), &parsed) {
            Ok(()) => Ok(true),
            Err(argon2::password_hash::Error::Password) => Ok(false),
            Err(e) => Err(AppError::Internal(anyhow::anyhow!("verify argon2: {e}"))),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::PasswordHasher;

    #[test]
    fn hash_and_verify_round_trip() {
        let h = PasswordHasher::new();
        let phc = h.hash("correct horse battery staple").unwrap();
        assert!(phc.starts_with("$argon2id$v=19$m=16384,t=4,p=2$"));
        assert!(h.verify("correct horse battery staple", &phc).unwrap());
        assert!(!h.verify("wrong password", &phc).unwrap());
    }
}
