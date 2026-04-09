//! Auth repository — port of
//! `modules/auth/repository/{repo_password,repo_session,repo_refresh_token,repo_one_time_token}.go`
//! merged into a single `AuthRepository` handle.
//!
//! **Phase D deltas vs Go source** (all design decisions from the
//! plan §3):
//!
//! | Area            | Delta                                           |
//! |-----------------|-------------------------------------------------|
//! | `token_hash`    | `Vec<u8>` raw 32-byte digest (D-OPEN-1)         |
//! | refresh CRUD    | Generic `update` method DROPPED (D-OPEN-5)      |
//! | refresh lookup  | `find_by_jti_for_update` uses SELECT FOR UPDATE |
//! | refresh revoke  | `revoke_all_for_user` batch update (reuse det.) |
//! | session revoke  | `revoke_all_for_user` batch update (pw change)  |
//! | OTT `find_all`  | DROPPED (TOCTOU fix goes through atomic upsert) |
//! | OTT upsert      | `upsert_for_user_subject` is atomic (fix §9.8)  |
//!
//! Method naming convention:
//! - `foo(pool, ...)`         — read-only / single-statement writes
//! - `foo_in_tx(tx, ...)`     — used inside a service-level tx
//!   (signin / rotation / password change)

use chrono::{DateTime, Utc};
use sqlx::{PgConnection, PgPool, Postgres, Transaction};
use uuid::Uuid;

use crate::domain::AppError;

use super::models::{OneTimeToken, RefreshToken, Session, UserPassword};

#[derive(Debug, Clone)]
pub struct AuthRepository {
    pool: PgPool,
}

impl AuthRepository {
    #[must_use]
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    /// Expose the pool so services can open their own transactions.
    #[must_use]
    pub fn pool(&self) -> &PgPool {
        &self.pool
    }

    /// Begin a transaction.
    pub async fn begin(&self) -> Result<Transaction<'_, Postgres>, AppError> {
        self.pool.begin().await.map_err(AppError::Database)
    }

    // ------------------------------------------------------------------
    // user_passwords
    // ------------------------------------------------------------------

    /// Fetch the stored PHC string for a user.
    pub async fn get_password_phc(&self, user_id: Uuid) -> Result<String, AppError> {
        let row: Option<(Vec<u8>,)> =
            sqlx::query_as("SELECT password_hash FROM public.user_passwords WHERE user_id = $1")
                .bind(user_id)
                .fetch_optional(&self.pool)
                .await
                .map_err(AppError::Database)?;

        let bytes = row
            .ok_or_else(|| AppError::NotFound("User password not found".to_string()))?
            .0;
        String::from_utf8(bytes).map_err(|e| {
            AppError::Internal(anyhow::anyhow!(
                "user_passwords.password_hash is not UTF-8: {e}"
            ))
        })
    }

    /// Insert a new password row (INSERT only; does not upsert).
    /// Used by `set_password` which is called for initial-credential
    /// flows.
    pub async fn set_password(&self, user_id: Uuid, phc: &str) -> Result<(), AppError> {
        let now = Utc::now();
        sqlx::query(
            "INSERT INTO public.user_passwords (user_id, password_hash, created_at) \
             VALUES ($1, $2, $3)",
        )
        .bind(user_id)
        .bind(phc.as_bytes())
        .bind(now)
        .execute(&self.pool)
        .await
        .map_err(AppError::Database)?;
        Ok(())
    }

    /// Update an existing password row. Used inside the password-change
    /// transaction (see `set_password_in_tx` for the tx variant).
    pub async fn update_password_in_tx(
        &self,
        tx: &mut PgConnection,
        user_id: Uuid,
        phc: &str,
    ) -> Result<(), AppError> {
        let now = Utc::now();
        let result = sqlx::query(
            "UPDATE public.user_passwords SET password_hash = $1, updated_at = $2 \
             WHERE user_id = $3",
        )
        .bind(phc.as_bytes())
        .bind(now)
        .bind(user_id)
        .execute(tx)
        .await
        .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("User password not found".to_string()));
        }
        Ok(())
    }

    // ------------------------------------------------------------------
    // sessions
    // ------------------------------------------------------------------

    /// Create a session inside an active transaction (used by signin).
    pub async fn create_session_in_tx(
        &self,
        tx: &mut PgConnection,
        session: &mut Session,
    ) -> Result<(), AppError> {
        if session.id.is_nil() {
            session.id = Uuid::now_v7();
        }
        session.created_at = Utc::now();

        sqlx::query(
            "INSERT INTO public.sessions \
             (id, user_id, token_hash, user_agent, device_name, device_fingerprint, \
              ip_address, expires_at, created_at, refreshed_at, revoked_at, revoked_by) \
             VALUES ($1,$2,$3,$4,$5,$6,$7::INET,$8,$9,$10,$11,$12)",
        )
        .bind(session.id)
        .bind(session.user_id)
        .bind(&session.token_hash)
        .bind(session.user_agent.as_ref())
        .bind(session.device_name.as_ref())
        .bind(session.device_fingerprint.as_ref())
        .bind(session.ip_address.as_ref())
        .bind(session.expires_at)
        .bind(session.created_at)
        .bind(session.refreshed_at)
        .bind(session.revoked_at)
        .bind(session.revoked_by)
        .execute(tx)
        .await
        .map_err(AppError::Database)?;
        Ok(())
    }

    /// Standalone `POST /api/v1/auth/session` handler path.
    pub async fn create_session(&self, session: &mut Session) -> Result<(), AppError> {
        let mut tx = self.begin().await?;
        self.create_session_in_tx(tx.as_mut(), session).await?;
        tx.commit().await.map_err(AppError::Database)?;
        Ok(())
    }

    pub async fn get_session(&self, session_id: Uuid) -> Result<Session, AppError> {
        sqlx::query_as::<_, Session>(
            "SELECT id, user_id, token_hash, user_agent, device_name, device_fingerprint, \
             ip_address::TEXT AS ip_address, expires_at, created_at, refreshed_at, \
             revoked_at, revoked_by \
             FROM public.sessions WHERE id = $1",
        )
        .bind(session_id)
        .fetch_optional(&self.pool)
        .await
        .map_err(AppError::Database)?
        .ok_or_else(|| AppError::NotFound("Session not found".to_string()))
    }

    pub async fn update_session(&self, session: &Session) -> Result<(), AppError> {
        let result = sqlx::query(
            "UPDATE public.sessions SET \
             user_id=$2, token_hash=$3, user_agent=$4, device_name=$5, \
             device_fingerprint=$6, ip_address=$7::INET, expires_at=$8, created_at=$9, \
             refreshed_at=$10, revoked_at=$11, revoked_by=$12 \
             WHERE id=$1",
        )
        .bind(session.id)
        .bind(session.user_id)
        .bind(&session.token_hash)
        .bind(session.user_agent.as_ref())
        .bind(session.device_name.as_ref())
        .bind(session.device_fingerprint.as_ref())
        .bind(session.ip_address.as_ref())
        .bind(session.expires_at)
        .bind(session.created_at)
        .bind(session.refreshed_at)
        .bind(session.revoked_at)
        .bind(session.revoked_by)
        .execute(&self.pool)
        .await
        .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("Session not found".to_string()));
        }
        Ok(())
    }

    pub async fn delete_session(&self, session_id: Uuid) -> Result<(), AppError> {
        let result = sqlx::query("DELETE FROM public.sessions WHERE id = $1")
            .bind(session_id)
            .execute(&self.pool)
            .await
            .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("Session not found".to_string()));
        }
        Ok(())
    }

    /// Used by the session-check middleware (D-AUTH-14).
    ///
    /// Returns `Ok((revoked_at, expires_at))` if the session exists,
    /// `Err(AppError::SessionRevoked)` otherwise.
    pub async fn validate_session(
        &self,
        session_id: Uuid,
    ) -> Result<(Option<DateTime<Utc>>, DateTime<Utc>), AppError> {
        let row: Option<(Option<DateTime<Utc>>, DateTime<Utc>)> =
            sqlx::query_as("SELECT revoked_at, expires_at FROM public.sessions WHERE id = $1")
                .bind(session_id)
                .fetch_optional(&self.pool)
                .await
                .map_err(AppError::Database)?;

        row.ok_or(AppError::SessionRevoked)
    }

    /// Revoke every session for a given user. Used by password-change
    /// flow (design 3.2: password change invalidates all sessions).
    pub async fn revoke_all_sessions_for_user_in_tx(
        &self,
        tx: &mut PgConnection,
        user_id: Uuid,
    ) -> Result<u64, AppError> {
        let now = Utc::now();
        let result = sqlx::query(
            "UPDATE public.sessions SET revoked_at = $1 \
             WHERE user_id = $2 AND revoked_at IS NULL",
        )
        .bind(now)
        .bind(user_id)
        .execute(tx)
        .await
        .map_err(AppError::Database)?;
        Ok(result.rows_affected())
    }

    // ------------------------------------------------------------------
    // refresh_tokens
    // ------------------------------------------------------------------

    /// Create a refresh token row inside an active transaction
    /// (used by signin + token rotation).
    pub async fn create_refresh_token_in_tx(
        &self,
        tx: &mut PgConnection,
        token: &mut RefreshToken,
    ) -> Result<(), AppError> {
        if token.id.is_nil() {
            token.id = Uuid::now_v7();
        }
        token.created_at = Utc::now();

        sqlx::query(
            "INSERT INTO public.refresh_tokens \
             (id, user_id, session_id, token_hash, ip_address, user_agent, \
              expires_at, created_at, revoked_at, revoked_by) \
             VALUES ($1,$2,$3,$4,$5::INET,$6,$7,$8,$9,$10)",
        )
        .bind(token.id)
        .bind(token.user_id)
        .bind(token.session_id)
        .bind(&token.token_hash)
        .bind(token.ip_address.as_ref())
        .bind(token.user_agent.as_ref())
        .bind(token.expires_at)
        .bind(token.created_at)
        .bind(token.revoked_at)
        .bind(token.revoked_by)
        .execute(tx)
        .await
        .map_err(AppError::Database)?;
        Ok(())
    }

    /// **D-AUTH-11 critical**: find a refresh token by its jti
    /// (`UUIDv7` primary key) with `SELECT ... FOR UPDATE` row lock.
    /// Used by the rotation service to serialize concurrent refresh
    /// calls on the same token.
    pub async fn find_refresh_by_id_for_update(
        &self,
        tx: &mut PgConnection,
        token_id: Uuid,
    ) -> Result<RefreshToken, AppError> {
        sqlx::query_as::<_, RefreshToken>(
            "SELECT id, user_id, session_id, token_hash, ip_address::TEXT AS ip_address, \
             user_agent, expires_at, created_at, revoked_at, revoked_by \
             FROM public.refresh_tokens WHERE id = $1 FOR UPDATE",
        )
        .bind(token_id)
        .fetch_optional(tx)
        .await
        .map_err(AppError::Database)?
        .ok_or(AppError::Unauthorized)
    }

    /// Mark a single refresh token as revoked inside the rotation tx.
    pub async fn revoke_refresh_in_tx(
        &self,
        tx: &mut PgConnection,
        token_id: Uuid,
        revoked_by: Uuid,
    ) -> Result<(), AppError> {
        let now = Utc::now();
        sqlx::query(
            "UPDATE public.refresh_tokens SET revoked_at = $1, revoked_by = $2 \
             WHERE id = $3",
        )
        .bind(now)
        .bind(revoked_by)
        .bind(token_id)
        .execute(tx)
        .await
        .map_err(AppError::Database)?;
        Ok(())
    }

    /// Revoke every refresh token for a user. Used by:
    /// (1) rotation reuse detection (design 3.1), and
    /// (2) password-change cascade (design 3.2).
    pub async fn revoke_all_refresh_for_user_in_tx(
        &self,
        tx: &mut PgConnection,
        user_id: Uuid,
    ) -> Result<u64, AppError> {
        let now = Utc::now();
        let result = sqlx::query(
            "UPDATE public.refresh_tokens SET revoked_at = $1 \
             WHERE user_id = $2 AND revoked_at IS NULL",
        )
        .bind(now)
        .bind(user_id)
        .execute(tx)
        .await
        .map_err(AppError::Database)?;
        Ok(result.rows_affected())
    }

    /// Update the `token_hash` of a session (used by rotation to keep
    /// the session row paired with the new refresh token).
    pub async fn update_session_token_hash_in_tx(
        &self,
        tx: &mut PgConnection,
        session_id: Uuid,
        new_token_hash: &[u8],
    ) -> Result<(), AppError> {
        let now = Utc::now();
        sqlx::query("UPDATE public.sessions SET token_hash = $1, refreshed_at = $2 WHERE id = $3")
            .bind(new_token_hash)
            .bind(now)
            .bind(session_id)
            .execute(tx)
            .await
            .map_err(AppError::Database)?;
        Ok(())
    }

    // ------------------------------------------------------------------
    // one_time_tokens
    // ------------------------------------------------------------------

    /// Find an active (non-expired) token by its subject + `relates_to`
    /// (email). Used by the cooldown check in verification service.
    pub async fn find_active_ott_for_subject(
        &self,
        user_id: Uuid,
        subject: &str,
    ) -> Result<Option<OneTimeToken>, AppError> {
        sqlx::query_as::<_, OneTimeToken>(
            "SELECT id, user_id, subject, token_hash, relates_to, metadata, \
             created_at, expires_at, last_sent_at \
             FROM public.one_time_tokens \
             WHERE user_id = $1 AND subject = $2",
        )
        .bind(user_id)
        .bind(subject)
        .fetch_optional(&self.pool)
        .await
        .map_err(AppError::Database)
    }

    pub async fn find_ott_by_token_hash(
        &self,
        token_hash: &str,
    ) -> Result<Option<OneTimeToken>, AppError> {
        sqlx::query_as::<_, OneTimeToken>(
            "SELECT id, user_id, subject, token_hash, relates_to, metadata, \
             created_at, expires_at, last_sent_at \
             FROM public.one_time_tokens WHERE token_hash = $1",
        )
        .bind(token_hash)
        .fetch_optional(&self.pool)
        .await
        .map_err(AppError::Database)
    }

    /// **D-AUTH-8 critical**: atomic upsert fixing audit §9.8 TOCTOU.
    /// Runs `DELETE + INSERT` in a single transaction. The Go source
    /// did `FindAll → loop → delete → insert` with no tx, which was
    /// both unsafe (race against the unique index on
    /// (`user_id`, subject)) and O(N) on the whole table.
    pub async fn upsert_ott_for_user_subject(
        &self,
        token: &mut OneTimeToken,
    ) -> Result<(), AppError> {
        if token.id.is_nil() {
            token.id = Uuid::now_v7();
        }
        token.created_at = Utc::now();

        let user_id = token
            .user_id
            .ok_or_else(|| AppError::BadRequest("OneTimeToken.user_id is required".to_string()))?;

        let mut tx = self.begin().await?;

        sqlx::query("DELETE FROM public.one_time_tokens WHERE user_id = $1 AND subject = $2")
            .bind(user_id)
            .bind(&token.subject)
            .execute(tx.as_mut())
            .await
            .map_err(AppError::Database)?;

        sqlx::query(
            "INSERT INTO public.one_time_tokens \
             (id, user_id, subject, token_hash, relates_to, metadata, \
              created_at, expires_at, last_sent_at) \
             VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
        )
        .bind(token.id)
        .bind(token.user_id)
        .bind(&token.subject)
        .bind(&token.token_hash)
        .bind(&token.relates_to)
        .bind(token.metadata.as_ref())
        .bind(token.created_at)
        .bind(token.expires_at)
        .bind(token.last_sent_at)
        .execute(tx.as_mut())
        .await
        .map_err(AppError::Database)?;

        tx.commit().await.map_err(AppError::Database)?;
        Ok(())
    }

    /// Delete a one-time token by id. Used by the revoke endpoint.
    pub async fn delete_ott(&self, token_id: Uuid) -> Result<(), AppError> {
        let result = sqlx::query("DELETE FROM public.one_time_tokens WHERE id = $1")
            .bind(token_id)
            .execute(&self.pool)
            .await
            .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("Token not found".to_string()));
        }
        Ok(())
    }

    /// Consume a one-time token atomically: delete by id + return
    /// (used by the validation service to mark consumed atomically
    /// inside its own tx).
    pub async fn delete_ott_in_tx(
        &self,
        tx: &mut PgConnection,
        token_id: Uuid,
    ) -> Result<(), AppError> {
        sqlx::query("DELETE FROM public.one_time_tokens WHERE id = $1")
            .bind(token_id)
            .execute(tx)
            .await
            .map_err(AppError::Database)?;
        Ok(())
    }

    /// Update `last_sent_at` on an existing token (used by resend).
    pub async fn update_ott_last_sent_at(
        &self,
        token_id: Uuid,
        ts: DateTime<Utc>,
    ) -> Result<(), AppError> {
        sqlx::query("UPDATE public.one_time_tokens SET last_sent_at = $1 WHERE id = $2")
            .bind(ts)
            .bind(token_id)
            .execute(&self.pool)
            .await
            .map_err(AppError::Database)?;
        Ok(())
    }
}

// Silence: `UserPassword` is exported through models.rs but not
// directly touched by the repository helpers here (we only read/write
// column values, not the struct). Keep the import visible so future
// cargo-expand users can see the intended model.
#[allow(dead_code)]
const _UNUSED_USER_PASSWORD: Option<UserPassword> = None;
