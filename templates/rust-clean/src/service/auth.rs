//! Auth service. Mirrors `apps/{{ package_name | kebab_case }}/service/auth.go`.

use std::sync::Arc;

use crate::config::env::Env;
use crate::domain::auth::LoginResponse;
use crate::domain::error::AppError;
use crate::repository::postgres::auth::AuthRepo;
use crate::utils::jwt::generate_token_pair;

pub struct AuthService {
    repo: Arc<AuthRepo>,
    env: Arc<Env>,
}

impl AuthService {
    pub fn new(repo: Arc<AuthRepo>, env: Arc<Env>) -> Self {
        Self { repo, env }
    }

    pub async fn login(&self, email: &str, password: &str) -> Result<LoginResponse, AppError> {
        let user = self.repo.authenticate_user(email, password).await?;
        let (access_token, refresh_token) = generate_token_pair(&user.id, &user.email, &self.env)?;
        Ok(LoginResponse {
            user,
            token: access_token,
            refresh_token,
        })
    }
}
