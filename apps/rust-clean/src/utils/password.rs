//! bcrypt password hashing. Mirrors `apps/go-clean/utils/bcrypt.go`.
//!
//! Uses `DEFAULT_COST` (matching Go's `bcrypt.DefaultCost`, currently 12).
//! `hash` and `verify` are synchronous + CPU-bound, so handlers wrap
//! them in `tokio::task::spawn_blocking` to avoid blocking the runtime.

use bcrypt::{DEFAULT_COST, hash, verify};

use crate::domain::error::AppError;

pub fn hash_password(password: &str) -> Result<String, AppError> {
    Ok(hash(password, DEFAULT_COST)?)
}

pub fn compare_password(password: &str, hash_str: &str) -> bool {
    verify(password, hash_str).unwrap_or(false)
}
