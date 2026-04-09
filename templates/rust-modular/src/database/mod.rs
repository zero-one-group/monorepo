//! Database pool / adapter — port of `internal/adapter/postgres.go`.
//!
//! Builds a sqlx `PgPool` with the pool-size and retry semantics from
//! `DatabaseConfig`. Retry loop matches the Go source's `PgMaxRetries`
//! behavior:
//! - `0`: no retries, fail immediately
//! - `N > 0`: retry N times with exponential backoff (200ms, 400ms, ...)
//! - `-1`: retry forever (used by integration tests that race DB boot)

use std::time::Duration;

use anyhow::{Context, Result, bail};
use sqlx::PgPool;
use sqlx::postgres::PgPoolOptions;
use tracing::{info, warn};

use crate::config::DatabaseConfig;

/// Open a Postgres connection pool with retry.
pub async fn connect_pool(cfg: &DatabaseConfig) -> Result<PgPool> {
    let max_retries = cfg.pg_max_retries;
    let infinite = max_retries < 0;
    // Use i32 throughout so the retries-remaining subtraction can't wrap.
    let max_retries_signed = max_retries.max(0);
    let mut attempt: i32 = 0;

    loop {
        attempt = attempt.saturating_add(1);
        match try_connect(cfg).await {
            Ok(pool) => {
                info!(
                    pg_max_pool_size = cfg.pg_max_pool_size,
                    attempt, "postgres pool connected"
                );
                return Ok(pool);
            }
            Err(err) => {
                let remaining = if infinite {
                    i32::MAX
                } else {
                    max_retries_signed.saturating_sub(attempt)
                };
                if !infinite && remaining <= 0 {
                    bail!("postgres pool connect failed after {attempt} attempts: {err:#}");
                }
                let backoff_ms: u64 = u64::from(attempt.min(5).unsigned_abs()) * 200;
                warn!(
                    attempt,
                    remaining = if infinite { -1 } else { remaining },
                    backoff_ms,
                    error = %err,
                    "postgres connect failed; retrying",
                );
                tokio::time::sleep(Duration::from_millis(backoff_ms)).await;
            }
        }
    }
}

async fn try_connect(cfg: &DatabaseConfig) -> Result<PgPool> {
    PgPoolOptions::new()
        .max_connections(cfg.pg_max_pool_size)
        .acquire_timeout(Duration::from_secs(10))
        .test_before_acquire(true)
        .connect(&cfg.database_url)
        .await
        .context("sqlx PgPoolOptions::connect")
}
