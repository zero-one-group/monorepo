//! Auth repository. Mirrors `apps/go-clean/internal/repository/postgres/auth.go`.
//!
//! `authenticate_user` looks up a user by email, verifies the password
//! against the bcrypt hash, and returns the (id, name, email) on success
//! or `AppError::Unauthorized` on any failure (mirrors the Go behavior
//! of collapsing "not found" and "wrong password" into the same error
//! to avoid user enumeration).

use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::error::AppError;
use crate::domain::user::User;
use crate::utils::password::compare_password;

pub struct AuthRepo {
    pool: PgPool,
}

impl AuthRepo {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    pub async fn authenticate_user(&self, email: &str, password: &str) -> Result<User, AppError> {
        let row: Option<(Uuid, String, String, String)> = sqlx::query_as(
            r"
            SELECT id, name, email, password
            FROM users
            WHERE email = $1 AND deleted_at IS NULL
            ",
        )
        .bind(email)
        .fetch_optional(&self.pool)
        .await?;
        let Some((id, name, email_db, hashed)) = row else {
            return Err(AppError::Unauthorized);
        };
        if !compare_password(password, &hashed) {
            return Err(AppError::Unauthorized);
        }
        Ok(User {
            id: id.to_string(),
            name,
            email: email_db,
            // The Go `authenticate_user` returns a User without populating
            // created_at/updated_at (zero values). We mirror that with
            // default timestamps; `login` only uses id + name + email.
            created_at: chrono::DateTime::<chrono::Utc>::from_timestamp(0, 0).unwrap(),
            updated_at: chrono::DateTime::<chrono::Utc>::from_timestamp(0, 0).unwrap(),
        })
    }
}
