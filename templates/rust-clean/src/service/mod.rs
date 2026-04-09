//! Service layer — orchestrates repositories and adds domain-level
//! logic on top. Mirrors `apps/{{ package_name | kebab_case }}/service/`.

pub mod auth;
pub mod user;
