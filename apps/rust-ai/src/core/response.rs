//! Response envelope types.
//!
//! Mirrors `app/core/response.py` so the JSON wire format is byte-equal:
//!
//! Success:
//! ```json
//! {
//!   "success": true,
//!   "message": "Operation successful",
//!   "data": { ... },
//!   "metadata": null,
//!   "error_code": null
//! }
//! ```
//!
//! Error (rendered by the `AppError` `IntoResponse` impl):
//! ```json
//! {
//!   "success": false,
//!   "message": "...",
//!   "error_code": "...",
//!   "data": null
//! }
//! ```
//!
//! NOTE on `exclude_none`: the Python `success_response()` calls
//! `.dict(exclude_none=True)`, which drops the `data`/`metadata`/`error_code`
//! keys when they're `None`. We replicate that with
//! `#[serde(skip_serializing_if = "Option::is_none")]` on every Optional
//! field of `SuccessResponse`.

use axum::Json;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use serde::Serialize;
use serde_json::{Value, json};

use crate::core::exception::AppError;

#[derive(Debug, Serialize)]
pub struct SuccessResponse<T: Serialize> {
    pub success: bool,
    pub message: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data: Option<T>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub metadata: Option<Value>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error_code: Option<String>,
}

impl<T: Serialize> SuccessResponse<T> {
    pub fn new(data: T) -> Self {
        Self {
            success: true,
            message: "Operation successful".to_string(),
            data: Some(data),
            metadata: None,
            error_code: None,
        }
    }
}

/// Mirror of Python `success_response(data, message="Operation successful", status_code=200)`.
pub fn success_response<T: Serialize>(data: T) -> Response {
    let body = SuccessResponse::new(data);
    (StatusCode::OK, Json(body)).into_response()
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let body = json!({
            "success": false,
            "message": self.message,
            "error_code": self.code,
            "data": self.data,
        });
        (self.status_code, Json(body)).into_response()
    }
}
