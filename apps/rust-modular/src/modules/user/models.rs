//! User models — port of `modules/user/models/{model,schema}.go`.
//!
//! Note: go-modular's Go `models/schema.go` only defines
//! `UserCreateRequest` (name + email). `UpdateUser` reuses the same
//! request shape in the Go handler (see `handler_update_user`). We
//! match that faithfully.

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::types::Json;
use utoipa::{IntoParams, ToSchema};
use uuid::Uuid;
use validator::Validate;

/// `UserMetadata` — the JSONB `metadata` column on the users table.
///
/// Only `timezone` is used by the current Go codebase (set to "UTC"
/// by default in `CreateUser`). Add fields here if go-modular grows
/// new metadata keys.
#[derive(Debug, Clone, Serialize, Deserialize, Default, ToSchema)]
pub struct UserMetadata {
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub timezone: Option<String>,
}

/// `User` — mirrors `modules/user/models/model.go` 1:1. All 13
/// columns from the `users` table. JSONB metadata is wrapped in
/// `sqlx::types::Json` so sqlx automatically decodes it as a struct
/// and serde serializes it inline (not as a double-encoded string).
#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow, ToSchema)]
pub struct User {
    pub id: Uuid,
    pub display_name: String,
    pub email: String,
    pub username: Option<String>,
    pub avatar_url: Option<String>,
    /// JSONB `metadata` column. Serialized inline as
    /// [`UserMetadata`] in both API responses and `OpenAPI`.
    #[schema(value_type = Option<UserMetadata>)]
    pub metadata: Option<Json<UserMetadata>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: Option<DateTime<Utc>>,
    pub email_verified_at: Option<DateTime<Utc>>,
    pub last_login_at: Option<DateTime<Utc>>,
    pub banned_at: Option<DateTime<Utc>>,
    pub ban_expires: Option<DateTime<Utc>>,
    pub ban_reason: Option<String>,
}

/// `UserCreateRequest` — JSON body for `POST /api/v1/users` and
/// `PUT /api/v1/users/:userId` (the Go handler reuses the same
/// request shape for create and update).
///
/// `username` is system-generated (see `UserService::create_user`),
/// not accepted from the client — matches the `json:"-"` tag in Go.
#[derive(Debug, Clone, Deserialize, Validate, ToSchema)]
pub struct UserCreateRequest {
    #[validate(length(min = 1, message = "name is required"))]
    #[schema(example = "Jane Doe")]
    pub name: String,

    #[validate(email(message = "invalid email format"))]
    #[schema(example = "jane@example.com")]
    pub email: String,
}

/// `FilterUser` — query parameters for `GET /api/v1/users`.
///
/// Matches `modules/user/models/model.go` `FilterUser`. The Go source
/// uses `query:"search"` / `query:"limit"` / `query:"offset"`; serde
/// reads the same parameter names via axum's `Query<T>` extractor.
#[derive(Debug, Clone, Deserialize, Default, IntoParams)]
#[into_params(parameter_in = Query)]
pub struct FilterUser {
    #[serde(default)]
    pub search: Option<String>,
    #[serde(default)]
    pub limit: Option<i64>,
    #[serde(default)]
    pub offset: Option<i64>,
}
