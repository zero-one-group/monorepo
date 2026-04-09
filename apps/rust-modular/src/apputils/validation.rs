//! Validation error formatter — port of
//! `pkg/apputils/validation.go::ValidationErrorsToMap`.
//!
//! Maps `validator::ValidationErrors` to a JSON object of the form
//! `{<field>: <message>}` with per-tag messages matching the Go
//! source strings byte-for-byte where possible.

use serde_json::{Map, Value};

/// Convert a `validator` error set into a JSON object suitable for
/// embedding under an `AppError::Validation` envelope.
///
/// The Go source uses the `json` struct tag to rename the field in
/// the output; `validator 0.20.0` already gives us the struct field
/// name via `field_errors()` keyed by the field identifier. For
/// complete byte-for-byte parity, callers should either ensure
/// their struct field names match the serde `rename` (we use the
/// `snake_case` default) or pass fields through `#[serde(rename)]`.
#[must_use]
pub fn validation_errors_to_map(errors: &validator::ValidationErrors) -> Value {
    let mut out = Map::new();
    for (field, field_errors) in errors.field_errors() {
        let Some(first) = field_errors.first() else {
            continue;
        };
        let msg = match first.code.as_ref() {
            "required" => format!("The {field} field is required"),
            "uuid" => "Must be a valid UUID".to_string(),
            "email" => "Must be a valid email address".to_string(),
            "length" => match first.params.get("min") {
                Some(min) => format!("Minimum length is {min}"),
                None => "Invalid length".to_string(),
            },
            "must_match" => match first.params.get("other") {
                Some(other) => format!("Must match {other}"),
                None => "Must match the related field".to_string(),
            },
            "url" => "Must be a valid URL".to_string(),
            _ => first
                .message
                .clone()
                .map_or_else(|| "Invalid value".to_string(), |m| m.to_string()),
        };
        out.insert((*field).to_string(), Value::String(msg));
    }
    Value::Object(out)
}
