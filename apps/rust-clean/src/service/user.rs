//! User service. Mirrors `apps/go-clean/service/user.go`.

use std::sync::Arc;

use uuid::Uuid;

use crate::domain::error::AppError;
use crate::domain::user::{CreateUserRequest, User, UserFilter};
use crate::repository::postgres::user::UserRepo;

pub struct UserService {
    repo: Arc<UserRepo>,
}

impl UserService {
    pub fn new(repo: Arc<UserRepo>) -> Self {
        Self { repo }
    }

    pub async fn create_user(&self, req: &CreateUserRequest) -> Result<User, AppError> {
        self.repo.create_user(req).await
    }

    pub async fn list_users(&self, filter: &UserFilter) -> Result<Vec<User>, AppError> {
        self.repo.list_users(filter).await
    }

    pub async fn get_user(&self, id: Uuid) -> Result<User, AppError> {
        match self.repo.get_user(id).await? {
            Some(user) => Ok(user),
            None => Err(AppError::UserNotFound),
        }
    }

    /// Update the name + email of an existing user. Mirrors the Go
    /// `UpdateUser` which first fetches the existing row, swaps the
    /// fields, and persists.
    pub async fn update_user(&self, id: Uuid, name: &str, email: &str) -> Result<User, AppError> {
        // Preserve the Go two-step pattern (fetch then update) so any
        // future business logic on "existing" stays in one place.
        let existing = self.repo.get_user(id).await?;
        if existing.is_none() {
            return Err(AppError::UserNotFound);
        }
        self.repo.update_user(id, name, email).await
    }

    pub async fn delete_user(&self, id: Uuid) -> Result<(), AppError> {
        let existing = self.repo.get_user(id).await?;
        if existing.is_none() {
            return Err(AppError::UserNotFound);
        }
        self.repo.delete_user(id).await
    }
}
