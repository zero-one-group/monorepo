//! Auth HTTP handlers — 14 endpoints.
//!
//! Handlers are thin passthroughs: they extract the request context
//! (IP, user-agent, device fingerprint) from axum headers and pass
//! it to the service layer where all Phase D logic lives.

use std::net::SocketAddr;

use axum::Json;
use axum::extract::{ConnectInfo, Extension, Path, Query, State};
use axum::http::{HeaderMap, StatusCode, header};
use axum::response::{IntoResponse, Redirect};
use serde_json::json;
use uuid::Uuid;
use validator::Validate;

use crate::AppState;
use crate::apputils::validation_errors_to_map;
use crate::domain::response::ErrorBody;
use crate::domain::{AppError, MessageResponse};

use super::middleware::AuthContext;
use super::models::{AuthenticatedUser, Session};
use super::schema::{
    CreateSessionRequest, InitiateEmailVerificationRequest, NeutralResponse,
    ResendEmailVerificationRequest, RevokeEmailVerificationRequest, SetPasswordRequest,
    SignInWithEmailRequest, SignInWithUsernameRequest, TokenRefreshRequest, UpdatePasswordRequest,
    UpdateSessionRequest, ValidateEmailVerificationRequest, VerifyEmailLinkQuery,
};
use super::service::RequestCtx;

// ----- helpers -----

fn ctx_from(headers: &HeaderMap, addr: &SocketAddr) -> RequestCtx {
    let user_agent = headers
        .get(header::USER_AGENT)
        .and_then(|h| h.to_str().ok())
        .map(ToString::to_string);
    let device_fingerprint = headers
        .get("x-device-fingerprint")
        .and_then(|h| h.to_str().ok())
        .map(ToString::to_string);
    let ip_address = headers
        .get("x-forwarded-for")
        .and_then(|h| h.to_str().ok())
        .and_then(|s| s.split(',').next())
        .map(str::trim)
        .map(ToString::to_string)
        .or_else(|| Some(addr.ip().to_string()));
    RequestCtx {
        ip_address,
        user_agent,
        device_fingerprint,
    }
}

fn parse_uuid(s: &str, field: &str) -> Result<Uuid, AppError> {
    Uuid::parse_str(s).map_err(|_| AppError::BadRequest(format!("{field} must be a valid UUID")))
}

fn validate<T: Validate>(req: &T) -> Result<(), AppError> {
    req.validate()
        .map_err(|e| AppError::Validation(validation_errors_to_map(&e)))
}

// ----- signin -----

#[utoipa::path(
    post,
    path = "/api/v1/auth/signin/email",
    tag = "Authentication",
    request_body = SignInWithEmailRequest,
    responses(
        (status = 200, description = "Signed in", body = AuthenticatedUser),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Invalid credentials or email not verified", body = ErrorBody),
    )
)]
pub async fn sign_in_with_email(
    State(state): State<AppState>,
    headers: HeaderMap,
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    Json(req): Json<SignInWithEmailRequest>,
) -> Result<Json<AuthenticatedUser>, AppError> {
    validate(&req)?;
    let ctx = ctx_from(&headers, &addr);
    let result = state.auth_service.sign_in_with_email(&req, &ctx).await?;
    Ok(Json(result))
}

#[utoipa::path(
    post,
    path = "/api/v1/auth/signin/username",
    tag = "Authentication",
    request_body = SignInWithUsernameRequest,
    responses(
        (status = 200, description = "Signed in", body = AuthenticatedUser),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Invalid credentials", body = ErrorBody),
    )
)]
pub async fn sign_in_with_username(
    State(state): State<AppState>,
    headers: HeaderMap,
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    Json(req): Json<SignInWithUsernameRequest>,
) -> Result<Json<AuthenticatedUser>, AppError> {
    validate(&req)?;
    let ctx = ctx_from(&headers, &addr);
    let result = state.auth_service.sign_in_with_username(&req, &ctx).await?;
    Ok(Json(result))
}

// ----- verify-email (link-based) -----

#[utoipa::path(
    get,
    path = "/api/v1/auth/verify-email",
    tag = "Email Verification",
    params(VerifyEmailLinkQuery),
    responses(
        (status = 200, description = "Email verified", body = MessageResponse),
        (status = 302, description = "Redirected to `redirect_to` URL"),
        (status = 401, description = "Invalid or expired token", body = ErrorBody),
    )
)]
pub async fn verify_email_by_link(
    State(state): State<AppState>,
    Query(q): Query<VerifyEmailLinkQuery>,
) -> impl IntoResponse {
    match state.auth_service.verify_email_by_link(&q.token).await {
        Ok(()) => {
            if let Some(redirect_url) = q.redirect_to {
                Redirect::to(&redirect_url).into_response()
            } else {
                Json(MessageResponse::new("Email verified")).into_response()
            }
        }
        Err(err) => err.into_response(),
    }
}

// ----- token rotation (NEW; replaces 4 deleted CRUD endpoints) -----

#[utoipa::path(
    post,
    path = "/api/v1/auth/token/refresh",
    tag = "Authentication",
    request_body = TokenRefreshRequest,
    responses(
        (status = 200, description = "Refreshed", body = AuthenticatedUser),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Invalid, expired, or reused refresh token", body = ErrorBody),
        (status = 409, description = "Concurrent refresh in progress", body = ErrorBody),
    )
)]
pub async fn rotate_refresh_token(
    State(state): State<AppState>,
    headers: HeaderMap,
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    Json(req): Json<TokenRefreshRequest>,
) -> Result<Json<AuthenticatedUser>, AppError> {
    validate(&req)?;
    let ctx = ctx_from(&headers, &addr);
    let result = state.auth_service.rotate_refresh_token(&req, &ctx).await?;
    Ok(Json(result))
}

// ----- verification (initiate, validate, revoke, resend) -----

#[utoipa::path(
    post,
    path = "/api/v1/auth/verification/email/initiate",
    tag = "Email Verification",
    request_body = InitiateEmailVerificationRequest,
    responses(
        (status = 202, description = "Verification email sent (neutral for unknown emails)", body = NeutralResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 429, description = "Cooldown — retry after Retry-After seconds", body = ErrorBody),
    )
)]
pub async fn initiate_email_verification(
    State(state): State<AppState>,
    Json(req): Json<InitiateEmailVerificationRequest>,
) -> Result<(StatusCode, Json<serde_json::Value>), AppError> {
    validate(&req)?;
    // Design 3.8: neutral response regardless of outcome (fix §9.25).
    // Errors propagate only for cooldown (429) and server errors.
    state.auth_service.initiate_email_verification(&req).await?;
    Ok((
        StatusCode::ACCEPTED,
        Json(json!({
            "message": "Verification email sent if the email is registered"
        })),
    ))
}

#[utoipa::path(
    post,
    path = "/api/v1/auth/verification/email/validate",
    tag = "Email Verification",
    request_body = ValidateEmailVerificationRequest,
    responses(
        (status = 200, description = "Email verified", body = MessageResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Invalid or expired token", body = ErrorBody),
    )
)]
pub async fn validate_email_verification(
    State(state): State<AppState>,
    Json(req): Json<ValidateEmailVerificationRequest>,
) -> Result<Json<MessageResponse>, AppError> {
    validate(&req)?;
    state.auth_service.validate_email_verification(&req).await?;
    Ok(Json(MessageResponse::new("Email verified successfully")))
}

#[utoipa::path(
    post,
    path = "/api/v1/auth/verification/email/revoke",
    tag = "Email Verification",
    request_body = RevokeEmailVerificationRequest,
    responses(
        (status = 200, description = "Verification token revoked", body = MessageResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn revoke_email_verification(
    State(state): State<AppState>,
    Json(req): Json<RevokeEmailVerificationRequest>,
) -> Result<Json<MessageResponse>, AppError> {
    validate(&req)?;
    state
        .auth_service
        .revoke_email_verification(&req.token)
        .await?;
    Ok(Json(MessageResponse::new("Verification token revoked")))
}

#[utoipa::path(
    post,
    path = "/api/v1/auth/verification/email/resend",
    tag = "Email Verification",
    request_body = ResendEmailVerificationRequest,
    responses(
        (status = 202, description = "Verification email resent (neutral for unknown emails)", body = NeutralResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 429, description = "Cooldown — retry after Retry-After seconds", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn resend_email_verification(
    State(state): State<AppState>,
    Json(req): Json<ResendEmailVerificationRequest>,
) -> Result<(StatusCode, Json<serde_json::Value>), AppError> {
    validate(&req)?;
    state.auth_service.resend_email_verification(&req).await?;
    Ok((
        StatusCode::ACCEPTED,
        Json(json!({
            "message": "Verification email sent if the email is registered"
        })),
    ))
}

// ----- password (ownership-checked) -----

#[utoipa::path(
    post,
    path = "/api/v1/auth/password",
    tag = "Password",
    request_body = SetPasswordRequest,
    responses(
        (status = 200, description = "Password set", body = MessageResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 403, description = "Cannot set password for another user", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn set_password(
    State(state): State<AppState>,
    Extension(auth): Extension<AuthContext>,
    Json(req): Json<SetPasswordRequest>,
) -> Result<Json<MessageResponse>, AppError> {
    validate(&req)?;
    state.auth_service.set_password(auth.user_id, &req).await?;
    Ok(Json(MessageResponse::new("Password set successfully")))
}

#[utoipa::path(
    put,
    path = "/api/v1/auth/password/{userId}",
    tag = "Password",
    params(("userId" = String, Path, description = "Target user UUID")),
    request_body = UpdatePasswordRequest,
    responses(
        (status = 200, description = "Password updated (all sessions + refresh tokens revoked)", body = MessageResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Invalid credentials", body = ErrorBody),
        (status = 403, description = "Cannot modify another user's password", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn update_password(
    State(state): State<AppState>,
    Extension(auth): Extension<AuthContext>,
    Path(user_id): Path<String>,
    Json(req): Json<UpdatePasswordRequest>,
) -> Result<Json<MessageResponse>, AppError> {
    validate(&req)?;
    let target_id = parse_uuid(&user_id, "user_id")?;
    state
        .auth_service
        .update_password(auth.user_id, target_id, &req)
        .await?;
    Ok(Json(MessageResponse::new("Password updated successfully")))
}

// ----- session CRUD (keep-track) -----

#[utoipa::path(
    post,
    path = "/api/v1/auth/session",
    tag = "Sessions",
    request_body = CreateSessionRequest,
    responses(
        (status = 201, description = "Session created", body = Session),
        (status = 400, description = "Validation failed", body = ErrorBody),
    )
)]
pub async fn create_session(
    State(state): State<AppState>,
    Json(req): Json<CreateSessionRequest>,
) -> Result<(StatusCode, Json<Session>), AppError> {
    validate(&req)?;
    let user_id = parse_uuid(&req.user_id, "user_id")?;
    // Decode hex token hash (client must send 64-char hex).
    let token_hash = hex::decode(&req.token_hash)
        .map_err(|_| AppError::BadRequest("token_hash must be hex-encoded".to_string()))?;
    if token_hash.len() != 32 {
        return Err(AppError::BadRequest(
            "token_hash must decode to 32 bytes".to_string(),
        ));
    }
    let expires_at = chrono::DateTime::parse_from_rfc3339(&req.expires_at)
        .map_err(|_| AppError::BadRequest("expires_at must be RFC3339".to_string()))?
        .with_timezone(&chrono::Utc);

    let mut session = Session {
        id: Uuid::nil(),
        user_id,
        token_hash,
        user_agent: req.user_agent,
        device_name: req.device_name,
        device_fingerprint: req.device_fingerprint,
        ip_address: req.ip_address,
        expires_at,
        created_at: chrono::Utc::now(),
        refreshed_at: None,
        revoked_at: None,
        revoked_by: None,
    };

    state.auth_service.create_session(&mut session).await?;
    Ok((StatusCode::CREATED, Json(session)))
}

#[utoipa::path(
    get,
    path = "/api/v1/auth/session/{sessionId}",
    tag = "Sessions",
    params(("sessionId" = String, Path, description = "Session UUID")),
    responses(
        (status = 200, description = "Session found", body = Session),
        (status = 400, description = "Invalid UUID", body = ErrorBody),
        (status = 404, description = "Session not found", body = ErrorBody),
    )
)]
pub async fn get_session(
    State(state): State<AppState>,
    Path(session_id): Path<String>,
) -> Result<Json<Session>, AppError> {
    let id = parse_uuid(&session_id, "session_id")?;
    let session = state.auth_service.get_session(id).await?;
    Ok(Json(session))
}

#[utoipa::path(
    put,
    path = "/api/v1/auth/session",
    tag = "Sessions",
    request_body = UpdateSessionRequest,
    responses(
        (status = 200, description = "Session updated", body = MessageResponse),
        (status = 400, description = "Validation failed", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 404, description = "Session not found", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn update_session(
    State(state): State<AppState>,
    Json(req): Json<UpdateSessionRequest>,
) -> Result<Json<MessageResponse>, AppError> {
    validate(&req)?;
    let session_id = parse_uuid(&req.session_id, "session_id")?;
    let mut existing = state.auth_service.get_session(session_id).await?;

    if req.user_agent.is_some() {
        existing.user_agent = req.user_agent;
    }
    if req.device_name.is_some() {
        existing.device_name = req.device_name;
    }
    if req.device_fingerprint.is_some() {
        existing.device_fingerprint = req.device_fingerprint;
    }
    if req.ip_address.is_some() {
        existing.ip_address = req.ip_address;
    }
    if let Some(ts) = req.refreshed_at {
        existing.refreshed_at = Some(
            chrono::DateTime::parse_from_rfc3339(&ts)
                .map_err(|_| AppError::BadRequest("refreshed_at must be RFC3339".to_string()))?
                .with_timezone(&chrono::Utc),
        );
    }
    if let Some(ts) = req.revoked_at {
        existing.revoked_at = Some(
            chrono::DateTime::parse_from_rfc3339(&ts)
                .map_err(|_| AppError::BadRequest("revoked_at must be RFC3339".to_string()))?
                .with_timezone(&chrono::Utc),
        );
    }
    if let Some(rb) = req.revoked_by {
        existing.revoked_by = Some(parse_uuid(&rb, "revoked_by")?);
    }

    state.auth_service.update_session(&existing).await?;
    Ok(Json(MessageResponse::new("Session updated successfully")))
}

#[utoipa::path(
    delete,
    path = "/api/v1/auth/session/{sessionId}",
    tag = "Sessions",
    params(("sessionId" = String, Path, description = "Session UUID")),
    responses(
        (status = 200, description = "Session deleted (access tokens invalidated)", body = MessageResponse),
        (status = 400, description = "Invalid UUID", body = ErrorBody),
        (status = 401, description = "Unauthorized", body = ErrorBody),
        (status = 404, description = "Session not found", body = ErrorBody),
    ),
    security(("bearer_auth" = []))
)]
pub async fn delete_session(
    State(state): State<AppState>,
    Path(session_id): Path<String>,
) -> Result<Json<MessageResponse>, AppError> {
    let id = parse_uuid(&session_id, "session_id")?;
    state.auth_service.delete_session(id).await?;
    Ok(Json(MessageResponse::new("Session deleted successfully")))
}
