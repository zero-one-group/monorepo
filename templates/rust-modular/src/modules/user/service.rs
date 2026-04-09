//! User service — port of `modules/user/services/services.go`.
//!
//! Thin business-logic layer over [`UserRepository`]. The only
//! non-trivial logic lives in `create_user`, which:
//!
//! 1. Assigns a `UUIDv7` id if the caller didn't supply one.
//! 2. Sets a default `UserMetadata { timezone: "UTC" }` if absent.
//! 3. Auto-generates a username from the email local-part, applying
//!    the same sanitize-and-numeric-suffix loop as the Go source.

use std::sync::Arc;

use sqlx::types::Json;
use uuid::Uuid;

use crate::domain::AppError;

use super::models::{FilterUser, User, UserMetadata};
use super::repository::UserRepository;

/// User service.
#[derive(Debug, Clone)]
pub struct UserService {
    repo: Arc<UserRepository>,
}

impl UserService {
    pub fn new(repo: Arc<UserRepository>) -> Self {
        Self { repo }
    }

    /// Create a new user with username auto-generation from email.
    ///
    /// Mirrors the Go `UserService.CreateUser` logic step-by-step.
    pub async fn create_user(&self, user: &mut User) -> Result<(), AppError> {
        if user.id.is_nil() {
            user.id = Uuid::now_v7();
        }

        if user.metadata.is_none() {
            user.metadata = Some(Json(UserMetadata {
                timezone: Some("UTC".to_string()),
            }));
        }

        // Generate username from email if not provided.
        if !user.email.is_empty()
            && user
                .username
                .as_ref()
                .is_none_or(std::string::String::is_empty)
        {
            let username = self.generate_unique_username(&user.email).await?;
            user.username = Some(username);
        }

        self.repo.create_user(user).await
    }

    async fn generate_unique_username(&self, email: &str) -> Result<String, AppError> {
        let base = email.split('@').next().unwrap_or("user").to_lowercase();

        // Matches the Go regex `[^a-z0-9_]+` without the regex dep:
        // keep ASCII alphanumerics and underscore, drop everything else.
        let mut sanitized: String = base
            .chars()
            .filter(|c| c.is_ascii_alphanumeric() || *c == '_')
            .collect();
        if sanitized.is_empty() {
            sanitized = "user".to_string();
        }

        let mut candidate = sanitized.clone();
        let mut suffix = 1u32;
        loop {
            if !self.repo.username_exists(&candidate).await? {
                return Ok(candidate);
            }
            candidate = format!("{sanitized}_{suffix}");
            suffix = suffix.saturating_add(1);
            if suffix > 10_000 {
                return Err(AppError::Conflict(
                    "Could not generate unique username".to_string(),
                ));
            }
        }
    }

    pub async fn get_user_by_id(&self, id: Uuid) -> Result<User, AppError> {
        self.repo.get_user_by_id(id).await
    }

    pub async fn list_users(&self, filter: &FilterUser) -> Result<Vec<User>, AppError> {
        self.repo.list_users(filter).await
    }

    pub async fn update_user(&self, user: &mut User) -> Result<(), AppError> {
        self.repo.update_user(user).await
    }

    pub async fn delete_user(&self, id: Uuid) -> Result<(), AppError> {
        self.repo.delete_user(id).await
    }

    pub async fn get_user_by_email(&self, email: &str) -> Result<User, AppError> {
        self.repo.get_user_by_email(email).await
    }

    pub async fn get_user_by_username(&self, username: &str) -> Result<User, AppError> {
        self.repo.get_user_by_username(username).await
    }

    pub async fn mark_email_verified(&self, user_id: Uuid) -> Result<(), AppError> {
        self.repo.mark_email_verified(user_id).await
    }
}
