//! Response envelope types. Mirrors `apps/{{ package_name | kebab_case }}/domain/response.go`.
//!
//! Field ordering matches the Go struct exactly so serde produces
//! byte-identical JSON:
//!
//! - `Response { code, message }`
//! - `ResponseSingleData<T> { code, data, message }`
//! - `ResponseMultipleData<T> { code, data, message }`
//!
//! **Important**: Go's `ResponseMultipleData` serializes `data` as a
//! JSON array even when the underlying slice is empty (Go's
//! `[]domain.User` marshals to `[]`, not `null`). We preserve this by
//! using `Vec<T>` without `#[serde(skip_serializing_if = "Vec::is_empty")]`.

use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
pub struct Response {
    pub code: u16,
    pub message: String,
}

impl Response {
    pub fn new(code: u16, message: impl Into<String>) -> Self {
        Self {
            code,
            message: message.into(),
        }
    }
}

#[derive(Debug, Clone, Serialize)]
pub struct ResponseSingleData<T: Serialize> {
    pub code: u16,
    pub data: T,
    pub message: String,
}

impl<T: Serialize> ResponseSingleData<T> {
    pub fn new(code: u16, data: T, message: impl Into<String>) -> Self {
        Self {
            code,
            data,
            message: message.into(),
        }
    }
}

#[derive(Debug, Clone, Serialize)]
pub struct ResponseMultipleData<T: Serialize> {
    pub code: u16,
    pub data: Vec<T>,
    pub message: String,
}

impl<T: Serialize> ResponseMultipleData<T> {
    pub fn new(code: u16, data: Vec<T>, message: impl Into<String>) -> Self {
        Self {
            code,
            data,
            message: message.into(),
        }
    }
}

/// Marker type used as `T` in `ResponseSingleData<Empty>` for responses
/// with a missing payload. Serializes to `{}`.
#[derive(Debug, Clone, Serialize, Default)]
pub struct Empty {}
