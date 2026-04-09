//! User module — 5 CRUD endpoints (keep-track).
//!
//! Port of `apps/{{ package_name | kebab_case }}/modules/user/` from Go to Rust. The
//! module follows a strict repository -> service -> handler layering.

pub mod handler;
pub mod models;
pub mod module;
pub mod repository;
pub mod service;

pub use module::routes;
pub use repository::UserRepository;
pub use service::UserService;
