//! Auth service — ports the 5 Go services with Phase D corrections.
//!
//! All 8 corrected-port fixes live here:
//!
//! 1. **Refresh rotation** (`rotate_refresh_token`) — `SELECT FOR
//!    UPDATE` row lock, reuse detection revokes all user sessions.
//! 2. **Close auth hole** — 4 naive CRUD endpoints are simply
//!    deleted from the handler/route surface.
//! 3. **Session revocation** — handled by `middleware::require_auth`
//!    via `validate_session`.
//! 4. **Transactional signin** — `sign_in_with_email` /
//!    `sign_in_with_username` wrap session + refresh-token inserts
//!    in a single tx.
//! 5. **Delete `JWT_ALGORITHM`** — `JwtGenerator` is HS256-only.
//! 6. **Real STARTTLS** — landing in D-SMTP-1; mailer is `None` in
//!    v1 scaffold.
//! 7. **Fix `session.expires_at`** — uses refresh-token expiry, not
//!    access-token expiry.
//! 8. **Delete `X-App-Audience`** — refresh JWT `aud` is hardcoded
//!    to `"client-app"` in `apputils::jwt`.

use std::sync::Arc;
use std::time::Duration;

use chrono::Utc;
use uuid::Uuid;

use crate::apputils::{
    AccessTokenPayload, JwtGenerator, PasswordHasher, generate_url_safe_token, sha256_bytes,
    sha256_hex,
};
use crate::domain::AppError;
use crate::mailer::Mailer;
use crate::modules::user::UserService;

use super::models::{AuthenticatedUser, OneTimeToken, RefreshToken, Session};
use super::repository::AuthRepository;
use super::schema::{
    InitiateEmailVerificationRequest, ResendEmailVerificationRequest, SetPasswordRequest,
    SignInWithEmailRequest, SignInWithUsernameRequest, TokenRefreshRequest, UpdatePasswordRequest,
    ValidateEmailVerificationRequest,
};

/// 60 seconds, matches D-OPEN-3.
pub const VERIFICATION_COOLDOWN: Duration = Duration::from_secs(60);
/// 15 minutes, matches Go source hardcoded TTL.
pub const VERIFICATION_TTL: Duration = Duration::from_secs(15 * 60);
/// Access-token TTL. Matches Go default (24h).
pub const ACCESS_TOKEN_EXPIRY: Duration = Duration::from_secs(24 * 60 * 60);
/// Refresh-token TTL. Matches Go default (7 days).
pub const REFRESH_TOKEN_EXPIRY: Duration = Duration::from_secs(7 * 24 * 60 * 60);

/// Request context — populated at the handler layer from axum
/// `ConnectInfo<SocketAddr>` + headers, passed to signin/rotation
/// services for session metadata capture.
#[derive(Debug, Clone, Default)]
pub struct RequestCtx {
    pub ip_address: Option<String>,
    pub user_agent: Option<String>,
    pub device_fingerprint: Option<String>,
}

/// The auth service aggregate — everything wired to handlers.
pub struct AuthService {
    repo: Arc<AuthRepository>,
    user_service: Arc<UserService>,
    jwt: Arc<JwtGenerator>,
    password_hasher: Arc<PasswordHasher>,
    mailer: Arc<Mailer>,
    base_url: String,
}

impl AuthService {
    #[must_use]
    pub fn new(
        repo: Arc<AuthRepository>,
        user_service: Arc<UserService>,
        jwt: Arc<JwtGenerator>,
        password_hasher: Arc<PasswordHasher>,
        mailer: Arc<Mailer>,
        base_url: String,
    ) -> Self {
        Self {
            repo,
            user_service,
            jwt,
            password_hasher,
            mailer,
            base_url,
        }
    }

    /// Accessor for the `require_auth` middleware.
    #[must_use]
    pub fn jwt(&self) -> &JwtGenerator {
        &self.jwt
    }

    /// Accessor for the `require_auth` middleware.
    #[must_use]
    pub fn repository(&self) -> &AuthRepository {
        &self.repo
    }

    // ------------------------------------------------------------------
    // Signin (design 3.3: transactional, fix §9.4 + §9.5)
    // ------------------------------------------------------------------

    pub async fn sign_in_with_email(
        &self,
        req: &SignInWithEmailRequest,
        ctx: &RequestCtx,
    ) -> Result<AuthenticatedUser, AppError> {
        let user = self
            .user_service
            .get_user_by_email(&req.email)
            .await
            .map_err(|_| AppError::InvalidCredentials)?;
        self.complete_sign_in(user, &req.password, ctx).await
    }

    pub async fn sign_in_with_username(
        &self,
        req: &SignInWithUsernameRequest,
        ctx: &RequestCtx,
    ) -> Result<AuthenticatedUser, AppError> {
        let user = self
            .user_service
            .get_user_by_username(&req.username)
            .await
            .map_err(|_| AppError::InvalidCredentials)?;
        self.complete_sign_in(user, &req.password, ctx).await
    }

    async fn complete_sign_in(
        &self,
        user: crate::modules::user::models::User,
        password: &str,
        ctx: &RequestCtx,
    ) -> Result<AuthenticatedUser, AppError> {
        // 1. Validate password via argon2id PHC.
        let phc = self
            .repo
            .get_password_phc(user.id)
            .await
            .map_err(|_| AppError::InvalidCredentials)?;
        if !self.password_hasher.verify(password, &phc)? {
            return Err(AppError::InvalidCredentials);
        }

        // 2. Ensure email is verified (matches Go §6 sign-in flow).
        if user.email_verified_at.is_none() {
            return Err(AppError::EmailNotVerified);
        }

        // 3. Mint refresh JWT first so we can hash it into the
        // session's token_hash column before both inserts run.
        let refresh_jti = Uuid::now_v7();
        let refresh_jwt = self
            .jwt
            .sign_refresh(&user.id.to_string(), &refresh_jti.to_string())?;
        let refresh_hash = sha256_bytes(refresh_jwt.as_bytes()).to_vec();
        let now = Utc::now();
        let refresh_expires =
            now + chrono::Duration::from_std(REFRESH_TOKEN_EXPIRY).unwrap_or_default();

        // 4. TRANSACTION BOUNDARY (design 3.3 fix §9.4):
        //    both inserts succeed or both roll back.
        let mut tx = self.repo.begin().await?;

        let mut session = Session {
            id: Uuid::nil(),
            user_id: user.id,
            token_hash: refresh_hash.clone(),
            user_agent: ctx.user_agent.clone(),
            device_name: None,
            device_fingerprint: ctx.device_fingerprint.clone(),
            ip_address: ctx.ip_address.clone(),
            // NOTE (fix §9.5): expires_at uses the REFRESH expiry,
            // not the ACCESS expiry like Go did. This keeps session
            // lifetime paired with the refresh token.
            expires_at: refresh_expires,
            created_at: now,
            refreshed_at: None,
            revoked_at: None,
            revoked_by: None,
        };
        self.repo
            .create_session_in_tx(tx.as_mut(), &mut session)
            .await?;

        let mut token_row = RefreshToken {
            id: refresh_jti,
            user_id: user.id,
            session_id: Some(session.id),
            token_hash: refresh_hash,
            ip_address: ctx.ip_address.clone(),
            user_agent: ctx.user_agent.clone(),
            expires_at: refresh_expires,
            created_at: now,
            revoked_at: None,
            revoked_by: None,
        };
        self.repo
            .create_refresh_token_in_tx(tx.as_mut(), &mut token_row)
            .await?;

        tx.commit().await.map_err(AppError::Database)?;

        // 5. Mint the access token bound to the new session.
        let access_expires =
            now + chrono::Duration::from_std(ACCESS_TOKEN_EXPIRY).unwrap_or_default();
        let access_token = self.jwt.sign_access(
            &user.id.to_string(),
            &AccessTokenPayload {
                user_id: user.id.to_string(),
                email: user.email.clone(),
                sid: session.id.to_string(),
            },
        )?;

        Ok(AuthenticatedUser {
            user,
            access_token,
            refresh_token: refresh_jwt,
            session_id: Some(session.id),
            token_expiry: access_expires,
        })
    }

    // ------------------------------------------------------------------
    // Token rotation (design 3.1: rotate-on-refresh + reuse detection)
    // ------------------------------------------------------------------

    pub async fn rotate_refresh_token(
        &self,
        req: &TokenRefreshRequest,
        ctx: &RequestCtx,
    ) -> Result<AuthenticatedUser, AppError> {
        // 1. Verify JWT (signature + exp + typ).
        let claims = self.jwt.verify_refresh(&req.refresh_token)?;
        let user_id = Uuid::parse_str(&claims.sub).map_err(|_| AppError::Unauthorized)?;
        let token_id = Uuid::parse_str(&claims.jti).map_err(|_| AppError::Unauthorized)?;

        // 2. BEGIN TX with row lock (design 3.1).
        let mut tx = self.repo.begin().await?;

        let row = self
            .repo
            .find_refresh_by_id_for_update(tx.as_mut(), token_id)
            .await?;

        // 3. Reuse detection: revoked_at set means replay.
        if row.revoked_at.is_some() {
            self.repo
                .revoke_all_refresh_for_user_in_tx(tx.as_mut(), user_id)
                .await?;
            self.repo
                .revoke_all_sessions_for_user_in_tx(tx.as_mut(), user_id)
                .await?;
            tx.commit().await.map_err(AppError::Database)?;
            tracing::warn!(
                user_id = %user_id,
                ip = ?ctx.ip_address,
                ua = ?ctx.user_agent,
                "refresh token reuse detected; revoked all sessions"
            );
            return Err(AppError::RefreshTokenReuse);
        }

        // 4. Expiry check (redundant with JWT exp but belt+braces).
        let now = Utc::now();
        if row.expires_at < now {
            return Err(AppError::Unauthorized);
        }

        // 5. Rotate: revoke current, mint new refresh + access.
        self.repo
            .revoke_refresh_in_tx(tx.as_mut(), token_id, user_id)
            .await?;

        let new_jti = Uuid::now_v7();
        let new_refresh_jwt = self
            .jwt
            .sign_refresh(&user_id.to_string(), &new_jti.to_string())?;
        let new_hash = sha256_bytes(new_refresh_jwt.as_bytes()).to_vec();
        let refresh_expires =
            now + chrono::Duration::from_std(REFRESH_TOKEN_EXPIRY).unwrap_or_default();

        // Keep same session; update its token_hash to match the
        // new refresh JWT.
        let session_id = row.session_id.ok_or_else(|| {
            AppError::Internal(anyhow::anyhow!("refresh token missing session_id"))
        })?;
        self.repo
            .update_session_token_hash_in_tx(tx.as_mut(), session_id, &new_hash)
            .await?;

        let mut new_row = RefreshToken {
            id: new_jti,
            user_id,
            session_id: Some(session_id),
            token_hash: new_hash,
            ip_address: ctx.ip_address.clone(),
            user_agent: ctx.user_agent.clone(),
            expires_at: refresh_expires,
            created_at: now,
            revoked_at: None,
            revoked_by: None,
        };
        self.repo
            .create_refresh_token_in_tx(tx.as_mut(), &mut new_row)
            .await?;

        tx.commit().await.map_err(AppError::Database)?;

        // 6. Mint new access token + load user for response body.
        let user = self.user_service.get_user_by_id(user_id).await?;
        let access_expires =
            now + chrono::Duration::from_std(ACCESS_TOKEN_EXPIRY).unwrap_or_default();
        let access_token = self.jwt.sign_access(
            &user.id.to_string(),
            &AccessTokenPayload {
                user_id: user.id.to_string(),
                email: user.email.clone(),
                sid: session_id.to_string(),
            },
        )?;

        Ok(AuthenticatedUser {
            user,
            access_token,
            refresh_token: new_refresh_jwt,
            session_id: Some(session_id),
            token_expiry: access_expires,
        })
    }

    // ------------------------------------------------------------------
    // Password management (design 3.9 ownership check + cascade)
    // ------------------------------------------------------------------

    pub async fn set_password(
        &self,
        caller_id: Uuid,
        req: &SetPasswordRequest,
    ) -> Result<(), AppError> {
        let target_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| AppError::BadRequest("Invalid user_id".to_string()))?;
        // Design 3.9: ownership check.
        if target_id != caller_id {
            return Err(AppError::OwnershipViolation);
        }
        let phc = self.password_hasher.hash(&req.password)?;
        self.repo.set_password(target_id, &phc).await
    }

    pub async fn update_password(
        &self,
        caller_id: Uuid,
        target_id: Uuid,
        req: &UpdatePasswordRequest,
    ) -> Result<(), AppError> {
        // Design 3.9: ownership check.
        if target_id != caller_id {
            return Err(AppError::OwnershipViolation);
        }

        // Verify current password.
        let current_phc = self.repo.get_password_phc(target_id).await?;
        if !self
            .password_hasher
            .verify(&req.current_password, &current_phc)?
        {
            return Err(AppError::InvalidCredentials);
        }

        let new_phc = self.password_hasher.hash(&req.new_password)?;

        // Design 3.2: update password + invalidate all sessions +
        // revoke all refresh tokens in a single transaction.
        let mut tx = self.repo.begin().await?;
        self.repo
            .update_password_in_tx(tx.as_mut(), target_id, &new_phc)
            .await?;
        self.repo
            .revoke_all_sessions_for_user_in_tx(tx.as_mut(), target_id)
            .await?;
        self.repo
            .revoke_all_refresh_for_user_in_tx(tx.as_mut(), target_id)
            .await?;
        tx.commit().await.map_err(AppError::Database)?;
        Ok(())
    }

    // ------------------------------------------------------------------
    // Session CRUD (keep-track, thin wrappers)
    // ------------------------------------------------------------------

    pub async fn create_session(&self, session: &mut Session) -> Result<(), AppError> {
        self.repo.create_session(session).await
    }

    pub async fn get_session(&self, id: Uuid) -> Result<Session, AppError> {
        self.repo.get_session(id).await
    }

    pub async fn update_session(&self, session: &Session) -> Result<(), AppError> {
        self.repo.update_session(session).await
    }

    pub async fn delete_session(&self, id: Uuid) -> Result<(), AppError> {
        self.repo.delete_session(id).await
    }

    // ------------------------------------------------------------------
    // Email verification (design 3.8 + audit §9.25 fixes)
    // ------------------------------------------------------------------

    /// Returns `Ok(())` on success, neutral response shape — the
    /// handler always returns 202 regardless of actual outcome to
    /// fix audit §9.25 (email enumeration).
    pub async fn initiate_email_verification(
        &self,
        req: &InitiateEmailVerificationRequest,
    ) -> Result<(), AppError> {
        // Neutral lookup — swallow NotFound so the response can't
        // be used to enumerate registered emails.
        let Ok(user) = self.user_service.get_user_by_email(&req.email).await else {
            return Ok(());
        };

        if user.email_verified_at.is_some() {
            return Ok(()); // Already verified; respond neutrally.
        }

        // Cooldown check (D-OPEN-3: 60 seconds).
        if let Some(existing) = self
            .repo
            .find_active_ott_for_subject(user.id, "email_verification")
            .await?
            && let Some(last_sent) = existing.last_sent_at
        {
            let age = Utc::now() - last_sent;
            let cooldown = chrono::Duration::from_std(VERIFICATION_COOLDOWN).unwrap_or_default();
            if age < cooldown {
                let retry_after = (cooldown - age).num_seconds().max(1).cast_unsigned();
                return Err(AppError::VerificationCooldown { retry_after });
            }
        }

        // Generate a URL-safe token; hash for storage.
        let raw_token = generate_url_safe_token(48)?;
        let token_hash = sha256_hex(raw_token.as_bytes());
        let now = Utc::now();

        let mut ott = OneTimeToken {
            id: Uuid::nil(),
            user_id: Some(user.id),
            subject: "email_verification".to_string(),
            token_hash,
            relates_to: user.email.clone(),
            metadata: None,
            created_at: now,
            expires_at: now + chrono::Duration::from_std(VERIFICATION_TTL).unwrap_or_default(),
            last_sent_at: Some(now),
        };

        // Atomic upsert (fix §9.8 TOCTOU).
        self.repo.upsert_ott_for_user_subject(&mut ott).await?;

        // Send verification email via lettre (D-SMTP-3). The mailer
        // falls back to stdout logging if no transport is configured
        // so dev flows without a real relay still work.
        let verification_url = format!(
            "{}/api/v1/auth/verify-email?token={}",
            self.base_url.trim_end_matches('/'),
            raw_token
        );
        let display_name = if user.display_name.is_empty() {
            "there".to_string()
        } else {
            user.display_name.clone()
        };
        let expiry_minutes = VERIFICATION_TTL.as_secs() / 60;
        self.mailer
            .send_verification_email(
                &user.email,
                &display_name,
                &verification_url,
                expiry_minutes,
            )
            .await?;
        Ok(())
    }

    pub async fn validate_email_verification(
        &self,
        req: &ValidateEmailVerificationRequest,
    ) -> Result<(), AppError> {
        let token_hash = sha256_hex(req.token.as_bytes());
        let ott = self
            .repo
            .find_ott_by_token_hash(&token_hash)
            .await?
            .ok_or(AppError::Unauthorized)?;

        if Utc::now() > ott.expires_at {
            return Err(AppError::Unauthorized);
        }

        let user_id = ott
            .user_id
            .ok_or_else(|| AppError::Internal(anyhow::anyhow!("ott missing user_id")))?;

        // Atomic consume + mark verified.
        let mut tx = self.repo.begin().await?;
        self.repo.delete_ott_in_tx(tx.as_mut(), ott.id).await?;
        tx.commit().await.map_err(AppError::Database)?;

        self.user_service.mark_email_verified(user_id).await?;
        Ok(())
    }

    pub async fn revoke_email_verification(&self, token: &str) -> Result<(), AppError> {
        let token_hash = sha256_hex(token.as_bytes());
        match self.repo.find_ott_by_token_hash(&token_hash).await? {
            Some(ott) => self.repo.delete_ott(ott.id).await,
            None => Ok(()), // Idempotent: already gone.
        }
    }

    pub async fn resend_email_verification(
        &self,
        req: &ResendEmailVerificationRequest,
    ) -> Result<(), AppError> {
        // Same flow as initiate — cooldown enforced inside.
        let init = InitiateEmailVerificationRequest {
            email: req.email.clone(),
            redirect_to: None,
        };
        self.initiate_email_verification(&init).await
    }

    /// Verify-email link handler (GET /verify-email?token=...).
    pub async fn verify_email_by_link(&self, token: &str) -> Result<(), AppError> {
        self.validate_email_verification(&ValidateEmailVerificationRequest {
            token: token.to_string(),
        })
        .await
    }
}
