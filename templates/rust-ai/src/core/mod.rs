//! Cross-cutting application primitives: env, database, exception, response,
//! instrumentation, logging. Mirrors `app/core/` in the original Python tree.

pub mod database;
pub mod env;
pub mod exception;
pub mod instrumentation;
pub mod logging;
pub mod response;
