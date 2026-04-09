//! User repository — port of
//! `apps/{{ package_name | kebab_case }}/modules/user/repository/repository.go`.
//!
//! sqlx-backed CRUD against the 13-column `public.users` table.
//! UUIDs are generated app-side (`Uuid::now_v7()`) because Postgres
//! 16 does not ship a `uuidv7()` function — see the 00001 migration
//! comment and D-OPEN-4 resolution.
//!
//! The Go source wraps `CreateUser` and `UpdateUser` in explicit
//! transactions. We do the same — even though each call issues just
//! one statement, the tx boundary matches the Go behavior and makes
//! the code safer if additional statements are added later.

use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::AppError;

use super::models::{FilterUser, User};

/// sqlx repository handle — owns a clone of the pool.
#[derive(Debug, Clone)]
pub struct UserRepository {
    pool: PgPool,
}

impl UserRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    /// Insert a user row. Generates `id` via `Uuid::now_v7()` if
    /// the caller passed `Uuid::nil()`. Matches the Go behavior
    /// exactly — `created_at` is set server-side via Rust
    /// `Utc::now()` (Go uses `time.Now()`).
    pub async fn create_user(&self, user: &mut User) -> Result<(), AppError> {
        if user.id.is_nil() {
            user.id = Uuid::now_v7();
        }
        user.created_at = Utc::now();

        let mut tx = self.pool.begin().await.map_err(AppError::Database)?;

        sqlx::query(
            "
            INSERT INTO public.users (
                id, display_name, email, username, avatar_url, metadata,
                created_at, updated_at, email_verified_at,
                last_login_at, banned_at, ban_expires, ban_reason
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
            ",
        )
        .bind(user.id)
        .bind(&user.display_name)
        .bind(&user.email)
        .bind(user.username.as_ref())
        .bind(user.avatar_url.as_ref())
        .bind(user.metadata.as_ref())
        .bind(user.created_at)
        .bind(user.updated_at)
        .bind(user.email_verified_at)
        .bind(user.last_login_at)
        .bind(user.banned_at)
        .bind(user.ban_expires)
        .bind(user.ban_reason.as_ref())
        .execute(&mut *tx)
        .await
        .map_err(AppError::Database)?;

        tx.commit().await.map_err(AppError::Database)?;
        Ok(())
    }

    /// `GET /api/v1/users/:userId`.
    pub async fn get_user_by_id(&self, id: Uuid) -> Result<User, AppError> {
        sqlx::query_as::<_, User>(
            "
            SELECT id, display_name, email, username, avatar_url, metadata,
                   created_at, updated_at, email_verified_at,
                   last_login_at, banned_at, ban_expires, ban_reason
            FROM public.users
            WHERE id = $1
            ",
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await
        .map_err(AppError::Database)?
        .ok_or_else(|| AppError::NotFound("User not found".to_string()))
    }

    /// `GET /api/v1/users` with optional search + pagination.
    ///
    /// Builds the query dynamically like the Go source, but uses
    /// sqlx's `QueryBuilder` for parameter binding instead of
    /// string concatenation with `fmt.Sprintf`.
    pub async fn list_users(&self, filter: &FilterUser) -> Result<Vec<User>, AppError> {
        let mut builder = sqlx::QueryBuilder::new(
            "SELECT id, display_name, email, username, avatar_url, metadata, \
             created_at, updated_at, email_verified_at, \
             last_login_at, banned_at, ban_expires, ban_reason \
             FROM public.users",
        );

        if let Some(search) = filter.search.as_ref().filter(|s| !s.is_empty()) {
            let pattern = format!("%{search}%");
            builder.push(" WHERE (display_name ILIKE ");
            builder.push_bind(pattern.clone());
            builder.push(" OR username ILIKE ");
            builder.push_bind(pattern);
            builder.push(")");
        }

        builder.push(" ORDER BY created_at DESC");

        if let Some(limit) = filter.limit.filter(|&l| l > 0) {
            builder.push(" LIMIT ");
            builder.push_bind(limit);
        }
        if let Some(offset) = filter.offset.filter(|&o| o > 0) {
            builder.push(" OFFSET ");
            builder.push_bind(offset);
        }

        let users = builder
            .build_query_as::<User>()
            .fetch_all(&self.pool)
            .await
            .map_err(AppError::Database)?;

        Ok(users)
    }

    /// Update an existing user row. Sets `updated_at` server-side.
    /// Returns `NotFound` if `RowsAffected == 0`.
    pub async fn update_user(&self, user: &mut User) -> Result<(), AppError> {
        user.updated_at = Some(Utc::now());

        let mut tx = self.pool.begin().await.map_err(AppError::Database)?;

        let result = sqlx::query(
            "
            UPDATE public.users
            SET display_name = $1,
                email = $2,
                username = $3,
                avatar_url = $4,
                metadata = $5,
                updated_at = $6,
                email_verified_at = $7,
                last_login_at = $8,
                banned_at = $9,
                ban_expires = $10,
                ban_reason = $11
            WHERE id = $12
            ",
        )
        .bind(&user.display_name)
        .bind(&user.email)
        .bind(user.username.as_ref())
        .bind(user.avatar_url.as_ref())
        .bind(user.metadata.as_ref())
        .bind(user.updated_at)
        .bind(user.email_verified_at)
        .bind(user.last_login_at)
        .bind(user.banned_at)
        .bind(user.ban_expires)
        .bind(user.ban_reason.as_ref())
        .bind(user.id)
        .execute(&mut *tx)
        .await
        .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("User not found".to_string()));
        }

        tx.commit().await.map_err(AppError::Database)?;
        Ok(())
    }

    /// Delete a user row. Returns `NotFound` on zero rows affected.
    pub async fn delete_user(&self, id: Uuid) -> Result<(), AppError> {
        let result = sqlx::query("DELETE FROM public.users WHERE id = $1")
            .bind(id)
            .execute(&self.pool)
            .await
            .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("User not found".to_string()));
        }
        Ok(())
    }

    /// Case-insensitive username existence check.
    pub async fn username_exists(&self, username: &str) -> Result<bool, AppError> {
        let row: Option<(i32,)> =
            sqlx::query_as("SELECT 1 FROM public.users WHERE LOWER(username) = LOWER($1) LIMIT 1")
                .bind(username)
                .fetch_optional(&self.pool)
                .await
                .map_err(AppError::Database)?;
        Ok(row.is_some())
    }

    /// Case-insensitive email existence check.
    pub async fn email_exists(&self, email: &str) -> Result<bool, AppError> {
        let row: Option<(i32,)> =
            sqlx::query_as("SELECT 1 FROM public.users WHERE LOWER(email) = LOWER($1) LIMIT 1")
                .bind(email)
                .fetch_optional(&self.pool)
                .await
                .map_err(AppError::Database)?;
        Ok(row.is_some())
    }

    /// Used by auth signin flow.
    pub async fn get_user_by_email(&self, email: &str) -> Result<User, AppError> {
        sqlx::query_as::<_, User>(
            "
            SELECT id, display_name, email, username, avatar_url, metadata,
                   created_at, updated_at, email_verified_at,
                   last_login_at, banned_at, ban_expires, ban_reason
            FROM public.users
            WHERE LOWER(email) = LOWER($1)
            LIMIT 1
            ",
        )
        .bind(email)
        .fetch_optional(&self.pool)
        .await
        .map_err(AppError::Database)?
        .ok_or_else(|| AppError::NotFound("User not found".to_string()))
    }

    /// Used by auth signin flow.
    pub async fn get_user_by_username(&self, username: &str) -> Result<User, AppError> {
        sqlx::query_as::<_, User>(
            "
            SELECT id, display_name, email, username, avatar_url, metadata,
                   created_at, updated_at, email_verified_at,
                   last_login_at, banned_at, ban_expires, ban_reason
            FROM public.users
            WHERE LOWER(username) = LOWER($1)
            LIMIT 1
            ",
        )
        .bind(username)
        .fetch_optional(&self.pool)
        .await
        .map_err(AppError::Database)?
        .ok_or_else(|| AppError::NotFound("User not found".to_string()))
    }

    /// Used by the email verification flow (D-AUTH-8). Sets
    /// `email_verified_at = NOW()` atomically.
    pub async fn mark_email_verified(&self, user_id: Uuid) -> Result<(), AppError> {
        let now = Utc::now();
        let result = sqlx::query("UPDATE public.users SET email_verified_at = $1 WHERE id = $2")
            .bind(now)
            .bind(user_id)
            .execute(&self.pool)
            .await
            .map_err(AppError::Database)?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("User not found".to_string()));
        }
        Ok(())
    }
}
