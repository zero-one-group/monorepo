//! User handlers: 5 endpoints under `/api/v1/users`.
//!
//! Mirrors `apps/{{ package_name | kebab_case }}/internal/rest/user.go`. All 5 handlers are
//! protected by the `require_auth` middleware wired at the nested
//! router in `rest::build`.

use axum::Json;
use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use uuid::Uuid;

use crate::AppState;
use crate::domain::response::{Empty, ResponseMultipleData, ResponseSingleData};
use crate::domain::user::{CreateUserRequest, UpdateUserRequest, User, UserFilter};

/// `GET /api/v1/users?search=...`
pub async fn list_users(
    State(state): State<AppState>,
    Query(filter): Query<UserFilter>,
) -> Response {
    match state.user_service.list_users(&filter).await {
        Ok(users) => {
            let body = ResponseMultipleData::<User>::new(
                StatusCode::OK.as_u16(),
                users,
                "Successfully retrieve user list",
            );
            (StatusCode::OK, Json(body)).into_response()
        }
        Err(err) => {
            let msg = format!("Failed to list users: {err}");
            let body = ResponseMultipleData::<Empty>::new(
                StatusCode::INTERNAL_SERVER_ERROR.as_u16(),
                vec![],
                msg,
            );
            (StatusCode::INTERNAL_SERVER_ERROR, Json(body)).into_response()
        }
    }
}

/// `GET /api/v1/users/{id}`
pub async fn get_user(State(state): State<AppState>, Path(id_param): Path<String>) -> Response {
    let Ok(id) = Uuid::parse_str(&id_param) else {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "Invalid user ID format",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    };
    match state.user_service.get_user(id).await {
        Ok(user) => {
            let body = ResponseSingleData::<User>::new(
                StatusCode::OK.as_u16(),
                user,
                "Successfully retrieved user",
            );
            (StatusCode::OK, Json(body)).into_response()
        }
        Err(crate::domain::error::AppError::UserNotFound) => {
            let body = ResponseSingleData::new(
                StatusCode::NOT_FOUND.as_u16(),
                Empty::default(),
                "User not found",
            );
            (StatusCode::NOT_FOUND, Json(body)).into_response()
        }
        Err(err) => {
            let msg = format!("Failed to get user: {err}");
            let body = ResponseSingleData::new(
                StatusCode::INTERNAL_SERVER_ERROR.as_u16(),
                Empty::default(),
                msg,
            );
            (StatusCode::INTERNAL_SERVER_ERROR, Json(body)).into_response()
        }
    }
}

/// `POST /api/v1/users`
pub async fn create_user(State(state): State<AppState>, body: axum::body::Bytes) -> Response {
    let Ok(req) = serde_json::from_slice::<CreateUserRequest>(&body) else {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "Invalid request payload",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    };
    match state.user_service.create_user(&req).await {
        Ok(user) => {
            let body = ResponseSingleData::<User>::new(
                StatusCode::CREATED.as_u16(),
                user,
                "User successfully created",
            );
            (StatusCode::CREATED, Json(body)).into_response()
        }
        Err(err) => {
            let msg = format!("Failed to create user: {err}");
            let body = ResponseSingleData::new(
                StatusCode::INTERNAL_SERVER_ERROR.as_u16(),
                Empty::default(),
                msg,
            );
            (StatusCode::INTERNAL_SERVER_ERROR, Json(body)).into_response()
        }
    }
}

/// `PUT /api/v1/users/{id}`
pub async fn update_user(
    State(state): State<AppState>,
    Path(id_param): Path<String>,
    body: axum::body::Bytes,
) -> Response {
    let Ok(id) = Uuid::parse_str(&id_param) else {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "Invalid user ID format",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    };
    let Ok(req) = serde_json::from_slice::<UpdateUserRequest>(&body) else {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "Invalid request payload",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    };
    match state
        .user_service
        .update_user(id, &req.name, &req.email)
        .await
    {
        Ok(user) => {
            let body = ResponseSingleData::<User>::new(
                StatusCode::OK.as_u16(),
                user,
                "User successfully updated",
            );
            (StatusCode::OK, Json(body)).into_response()
        }
        Err(err) => {
            let msg = format!("Failed to update user: {err}");
            let body = ResponseSingleData::new(
                StatusCode::INTERNAL_SERVER_ERROR.as_u16(),
                Empty::default(),
                msg,
            );
            (StatusCode::INTERNAL_SERVER_ERROR, Json(body)).into_response()
        }
    }
}

/// `DELETE /api/v1/users/{id}`
pub async fn delete_user(State(state): State<AppState>, Path(id_param): Path<String>) -> Response {
    let Ok(id) = Uuid::parse_str(&id_param) else {
        let body = ResponseSingleData::new(
            StatusCode::BAD_REQUEST.as_u16(),
            Empty::default(),
            "Invalid user ID format",
        );
        return (StatusCode::BAD_REQUEST, Json(body)).into_response();
    };
    match state.user_service.delete_user(id).await {
        Ok(()) => {
            let body = ResponseSingleData::new(
                StatusCode::NO_CONTENT.as_u16(),
                Empty::default(),
                "User successfully deleted",
            );
            (StatusCode::NO_CONTENT, Json(body)).into_response()
        }
        Err(err) => {
            let msg = format!("Failed to delete user: {err}");
            let body = ResponseSingleData::new(
                StatusCode::INTERNAL_SERVER_ERROR.as_u16(),
                Empty::default(),
                msg,
            );
            (StatusCode::INTERNAL_SERVER_ERROR, Json(body)).into_response()
        }
    }
}
