//! Database connection pool and health check.
//!
//! Mirrors `app/core/database.py`. The Python service uses `SQLAlchemy`
//! async with a singleton pattern; we use an `sqlx::PgPool` wrapped in
//! `Arc` so it's cheaply cloneable into [`AppState`] and every handler.
//!
//! `check_connection()` runs `SELECT 1` and returns `bool`, matching the
//! Python `Database.check_connection()` semantics exactly.

use std::time::Duration;

use anyhow::{Context, Result};
use sqlx::PgPool;
use sqlx::postgres::PgPoolOptions;

#[derive(Clone)]
pub struct Database {
    pool: PgPool,
}

impl Database {
    pub async fn connect(database_url: &str) -> Result<Self> {
        let pool = PgPoolOptions::new()
            .max_connections(10)
            .acquire_timeout(Duration::from_secs(10))
            .test_before_acquire(true) // pool_pre_ping=True equivalent
            .connect(database_url)
            .await
            .context("create pg pool")?;
        tracing::debug!("Database initialized successfully");
        Ok(Self { pool })
    }

    /// Construct from an existing `PgPool`.
    ///
    /// Used by integration tests that get their pool from a
    /// testcontainer-managed Postgres instance instead of booting
    /// from a `DATABASE_URL`.
    #[must_use]
    pub fn from_pool(pool: PgPool) -> Self {
        Self { pool }
    }

    pub fn pool(&self) -> &PgPool {
        &self.pool
    }

    /// Run `SELECT 1` against the pool. Returns `false` on any error,
    /// matching the Python `check_connection()` truthiness contract.
    pub async fn check_connection(&self) -> bool {
        match sqlx::query("SELECT 1").execute(&self.pool).await {
            Ok(_) => {
                tracing::debug!("Database connection check successful.");
                true
            }
            Err(e) => {
                tracing::error!(error = %e, "Database connection check failed");
                false
            }
        }
    }

    /// Drain and close the pool. Called from the shutdown handler.
    pub async fn dispose(&self) {
        self.pool.close().await;
        tracing::info!("Database engine disposed");
    }
}
