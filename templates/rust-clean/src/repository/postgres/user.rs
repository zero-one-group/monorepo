//! User repository. Mirrors `apps/{{ package_name | kebab_case }}/internal/repository/postgres/user.go`.

use chrono::{DateTime, Utc};
use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::error::AppError;
use crate::domain::user::{CreateUserRequest, User, UserFilter};
use crate::utils::password::hash_password;

/// Row shape returned by every user SELECT query. Factored to keep
/// clippy's `type_complexity` lint quiet — rust-analyzer otherwise
/// flags the inline `(Uuid, String, String, DateTime<Utc>, DateTime<Utc>)`
/// tuple as too complex.
type UserRow = (Uuid, String, String, DateTime<Utc>, DateTime<Utc>);

fn user_from_row(row: UserRow) -> User {
    let (id, name, email, created_at, updated_at) = row;
    User {
        id: id.to_string(),
        name,
        email,
        created_at,
        updated_at,
    }
}

pub struct UserRepo {
    pool: PgPool,
}

impl UserRepo {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    /// Insert a new user. Returns the inserted row (id auto-generated
    /// by Postgres via `uuid_generate_v4()`).
    pub async fn create_user(&self, req: &CreateUserRequest) -> Result<User, AppError> {
        let id = Uuid::new_v4();
        let hashed = hash_password(&req.password)?;
        let row: (DateTime<Utc>, DateTime<Utc>) = sqlx::query_as(
            r"
            INSERT INTO users (id, name, email, password, created_at, updated_at)
            VALUES ($1, $2, $3, $4, NOW(), NOW())
            RETURNING created_at, updated_at
            ",
        )
        .bind(id)
        .bind(&req.name)
        .bind(&req.email)
        .bind(&hashed)
        .fetch_one(&self.pool)
        .await?;
        Ok(User {
            id: id.to_string(),
            name: req.name.clone(),
            email: req.email.clone(),
            created_at: row.0,
            updated_at: row.1,
        })
    }

    /// List all non-deleted users, optionally filtered by name/email LIKE.
    pub async fn list_users(&self, filter: &UserFilter) -> Result<Vec<User>, AppError> {
        let rows: Vec<UserRow> =
            if let Some(search) = filter.search.as_deref().filter(|s| !s.is_empty()) {
                let like = format!("%{search}%");
                sqlx::query_as(
                    r"
                    SELECT u.id, u.name, u.email, u.created_at, u.updated_at
                    FROM users u
                    WHERE u.deleted_at IS NULL
                      AND (u.name ILIKE $1 OR u.email ILIKE $1)
                    ",
                )
                .bind(like)
                .fetch_all(&self.pool)
                .await?
            } else {
                sqlx::query_as(
                    r"
                    SELECT u.id, u.name, u.email, u.created_at, u.updated_at
                    FROM users u
                    WHERE u.deleted_at IS NULL
                    ",
                )
                .fetch_all(&self.pool)
                .await?
            };
        Ok(rows.into_iter().map(user_from_row).collect())
    }

    /// Fetch one non-deleted user by ID. Returns `Ok(None)` if not found.
    pub async fn get_user(&self, id: Uuid) -> Result<Option<User>, AppError> {
        let row: Option<UserRow> = sqlx::query_as(
            r"
                SELECT id, name, email, created_at, updated_at
                FROM users
                WHERE id = $1 AND deleted_at IS NULL
                ",
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;
        Ok(row.map(user_from_row))
    }

    /// Update an existing user's name + email. Returns the updated row
    /// or `AppError::UserNotFound` if the row is gone or soft-deleted.
    pub async fn update_user(&self, id: Uuid, name: &str, email: &str) -> Result<User, AppError> {
        let row: Option<UserRow> = sqlx::query_as(
            r"
                UPDATE users
                SET name = $1,
                    email = $2,
                    updated_at = NOW()
                WHERE id = $3 AND deleted_at IS NULL
                RETURNING id, name, email, created_at, updated_at
                ",
        )
        .bind(name)
        .bind(email)
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;
        match row {
            Some(r) => Ok(user_from_row(r)),
            None => Err(AppError::UserNotFound),
        }
    }

    /// Soft-delete a user. Returns `AppError::UserNotFound` if no row
    /// was affected (row already deleted or never existed).
    pub async fn delete_user(&self, id: Uuid) -> Result<(), AppError> {
        let result = sqlx::query(
            r"
            UPDATE users
            SET deleted_at = NOW()
            WHERE id = $1 AND deleted_at IS NULL
            ",
        )
        .bind(id)
        .execute(&self.pool)
        .await?;
        if result.rows_affected() == 0 {
            return Err(AppError::UserNotFound);
        }
        Ok(())
    }
}
