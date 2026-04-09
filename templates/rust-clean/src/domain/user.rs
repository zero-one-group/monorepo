//! User domain types. Mirrors `apps/{{ package_name | kebab_case }}/domain/user.go`.

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: String,
    pub name: String,
    pub email: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Deserialize)]
pub struct CreateUserRequest {
    pub name: String,
    pub email: String,
    pub password: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct UpdateUserRequest {
    pub name: String,
    pub email: String,
}

#[derive(Debug, Clone, Default, Deserialize)]
pub struct UserFilter {
    #[serde(default)]
    pub search: Option<String>,
}
