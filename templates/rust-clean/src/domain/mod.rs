//! Domain types.
//!
//! Mirrors `apps/{{ package_name | kebab_case }}/domain/`:
//! - `user` → `User`, `CreateUserRequest`, `UpdateUserRequest`, `UserFilter`
//! - `auth` → `LoginRequest`, `LoginResponse`, `JwtClaim`
//! - `response` → `Response`, `ResponseSingleData<T>`, `ResponseMultipleData<T>`, `Empty`
//! - `error` → `AppError`

pub mod auth;
pub mod error;
pub mod response;
pub mod user;

pub use auth::{JwtClaim, LoginRequest, LoginResponse};
pub use error::AppError;
pub use response::{Empty, Response, ResponseMultipleData, ResponseSingleData};
pub use user::{CreateUserRequest, UpdateUserRequest, User, UserFilter};
