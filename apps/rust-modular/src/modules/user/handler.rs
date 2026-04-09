//! User HTTP handlers — port of
//! `modules/user/handler/handler.go`.
//!
//! 5 endpoints:
//!   - `POST   /api/v1/users`           `create_user`
//!   - `GET    /api/v1/users`           `list_users`
//!   - `GET    /api/v1/users/:userId`   `get_user`
//!   - `PUT    /api/v1/users/:userId`   `update_user`
//!   - `DELETE /api/v1/users/:userId`   `delete_user`
//!
//! Response shape follows the Go source convention:
//! - **Success (create/list/get)**: raw `User` / `Vec<User>` JSON
//! - **Success (update/delete)**:    `{"message": "..."}`
//! - **Errors**:                      `{"error": "..."}` + optional `{"details": ...}`
//!
//! Validation uses the `validator` crate. Per-field errors are
//! serialized as `{<field>: <message>}` into the `details` key,
//! matching the shape of Go's `apputils.ValidationErrorsToMap`.

use axum::Json;
use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use serde_json::{Map, Value};
use uuid::Uuid;
use validator::Validate;

use crate::AppState;
use crate::domain::response::ErrorBody;
use crate::domain::{AppError, MessageResponse};

use super::models::{FilterUser, User, UserCreateRequest};

/// `POST /api/v1/users`.
#[utoipa::path(
    post,
    path = "/api/v1/users",
    tag = "User Management",
    request_body = UserCreateRequest,
    responses(
        (status = 201, description = "User created", body = User),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn create_user(
    State(state): State<AppState>,
    Json(req): Json<UserCreateRequest>,
) -> Result<(StatusCode, Json<User>), AppError> {
    if let Err(errors) = req.validate() {
        return Err(AppError::Validation(validation_errors_to_details(&errors)));
    }

    let mut user = User {
        id: Uuid::nil(),
        display_name: req.name,
        email: req.email,
        username: None,
        avatar_url: None,
        metadata: None,
        created_at: chrono::Utc::now(),
        updated_at: None,
        email_verified_at: None,
        last_login_at: None,
        banned_at: None,
        ban_expires: None,
        ban_reason: None,
    };

    state.user_service.create_user(&mut user).await?;
    Ok((StatusCode::CREATED, Json(user)))
}

/// `GET /api/v1/users`.
#[utoipa::path(
    get,
    path = "/api/v1/users",
    tag = "User Management",
    params(FilterUser),
    responses(
        (status = 200, description = "List users", body = Vec<User>),
        (status = 401, description = "Unauthorized", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn list_users(
    State(state): State<AppState>,
    Query(filter): Query<FilterUser>,
) -> Result<Json<Vec<User>>, AppError> {
    let users = state.user_service.list_users(&filter).await?;
    Ok(Json(users))
}

/// `GET /api/v1/users/:userId`.
#[utoipa::path(
    get,
    path = "/api/v1/users/{userId}",
    tag = "User Management",
    params(("userId" = String, Path, description = "User UUID")),
    responses(
        (status = 200, description = "User found", body = User),
        (status = 400, description = "Invalid UUID", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 404, description = "User not found", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn get_user(
    State(state): State<AppState>,
    Path(user_id): Path<String>,
) -> Result<Json<User>, AppError> {
    let id = parse_user_id(&user_id)?;
    let user = state.user_service.get_user_by_id(id).await?;
    Ok(Json(user))
}

/// `PUT /api/v1/users/:userId`.
#[utoipa::path(
    put,
    path = "/api/v1/users/{userId}",
    tag = "User Management",
    params(("userId" = String, Path, description = "User UUID")),
    request_body = UserCreateRequest,
    responses(
        (status = 200, description = "User updated", body = MessageResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 404, description = "User not found", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn update_user(
    State(state): State<AppState>,
    Path(user_id): Path<String>,
    Json(req): Json<UserCreateRequest>,
) -> Result<Json<MessageResponse>, AppError> {
    let id = parse_user_id(&user_id)?;
    if let Err(errors) = req.validate() {
        return Err(AppError::Validation(validation_errors_to_details(&errors)));
    }

    let mut user = User {
        id,
        display_name: req.name,
        email: req.email,
        username: None,
        avatar_url: None,
        metadata: None,
        created_at: chrono::Utc::now(),
        updated_at: None,
        email_verified_at: None,
        last_login_at: None,
        banned_at: None,
        ban_expires: None,
        ban_reason: None,
    };

    state.user_service.update_user(&mut user).await?;
    Ok(Json(MessageResponse::new("User updated successfully")))
}

/// `DELETE /api/v1/users/:userId`.
#[utoipa::path(
    delete,
    path = "/api/v1/users/{userId}",
    tag = "User Management",
    params(("userId" = String, Path, description = "User UUID")),
    responses(
        (status = 200, description = "User deleted", body = MessageResponse),
        (status = 400, description = "Invalid UUID", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 404, description = "User not found", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn delete_user(
    State(state): State<AppState>,
    Path(user_id): Path<String>,
) -> Result<Json<MessageResponse>, AppError> {
    let id = parse_user_id(&user_id)?;
    state.user_service.delete_user(id).await?;
    Ok(Json(MessageResponse::new("User deleted successfully")))
}

fn parse_user_id(s: &str) -> Result<Uuid, AppError> {
    Uuid::parse_str(s)
        .map_err(|_| AppError::BadRequest("User ID in path must be a valid UUID".to_string()))
}

/// Port of Go's `apputils.ValidationErrorsToMap`. Produces a
/// `{"<field>": "<message>"}` JSON object from a `validator` error
/// set.
fn validation_errors_to_details(errors: &validator::ValidationErrors) -> Value {
    let mut details = Map::new();
    for (field, field_errors) in errors.field_errors() {
        if let Some(first) = field_errors.first() {
            let msg = first
                .message
                .clone()
                .map_or_else(|| first.code.to_string(), |m| m.to_string());
            details.insert((*field).to_string(), Value::String(msg));
        }
    }
    Value::Object(details)
}
