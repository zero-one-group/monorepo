//! Rust port of the Zero One Group go-clean service.
//!
//! Behavior-equivalent to the original Go `apps/go-clean/` tree:
//! - 7 HTTP endpoints (health + auth login + 5-endpoint user CRUD)
//! - Clean-architecture layers: `domain` → `service` → `repository` →
//!   `rest`. Kept as four top-level modules so the code is discoverable
//!   by anyone who previously worked on the Go version.
//! - HS256 JWT auth via `jsonwebtoken` (mirrors `golang-jwt/jwt/v5`)
//! - bcrypt password hashing via the `bcrypt` crate (mirrors
//!   `golang.org/x/crypto/bcrypt` at `DefaultCost`)
//! - Per-IP rate limiting (10 req/s, burst 20) via `tower_governor`
//! - Same env var contract (`APP_HOST`, `APP_PORT`, `DATABASE_URL`,
//!   `JWT_SECRET`, `AUTH_TOKEN_EXPIRY_MINUTES`, etc.)
//!
//! Public surface for tests and the binary entrypoint:
//! - [`serve`] — boots the HTTP server, blocks until shutdown
//! - [`build_router`] — constructs the [`axum::Router`] alone (used by
//!   integration tests so they don't bind a real port)
//! - [`AppState`] — shared state (DB pool + JWT config)

pub mod config;
pub mod domain;
pub mod repository;
pub mod rest;
pub mod service;
pub mod utils;

use std::net::SocketAddr;
use std::sync::Arc;
use std::time::Duration;

use anyhow::{Context, Result};
use sqlx::PgPool;
use sqlx::postgres::PgPoolOptions;
use tokio::net::TcpListener;
use tokio::signal;
use tracing::info;

use crate::config::env::Env;
use crate::repository::postgres::auth::AuthRepo;
use crate::repository::postgres::user::UserRepo;
use crate::service::auth::AuthService;
use crate::service::user::UserService;

/// Application state shared with every handler via `axum::extract::State`.
#[derive(Clone)]
pub struct AppState {
    pub env: Arc<Env>,
    pub pool: PgPool,
    pub user_service: Arc<UserService>,
    pub auth_service: Arc<AuthService>,
}

impl AppState {
    /// Construct from a loaded [`Env`]. Creates the DB pool and wires
    /// the repository + service layers up-front.
    pub async fn from_env(env: Env) -> Result<Self> {
        let pool = PgPoolOptions::new()
            .max_connections(10)
            .acquire_timeout(Duration::from_secs(10))
            .test_before_acquire(true)
            .connect(&env.database_url)
            .await
            .context("create pg pool")?;

        let env = Arc::new(env);
        let user_repo = Arc::new(UserRepo::new(pool.clone()));
        let auth_repo = Arc::new(AuthRepo::new(pool.clone()));
        let user_service = Arc::new(UserService::new(user_repo));
        let auth_service = Arc::new(AuthService::new(auth_repo, env.clone()));

        Ok(Self {
            env,
            pool,
            user_service,
            auth_service,
        })
    }

    /// Test-only constructor: inject an existing `PgPool` (from a
    /// testcontainer) alongside a pre-built `Env`.
    pub fn from_parts(env: Env, pool: PgPool) -> Self {
        let env = Arc::new(env);
        let user_repo = Arc::new(UserRepo::new(pool.clone()));
        let auth_repo = Arc::new(AuthRepo::new(pool.clone()));
        let user_service = Arc::new(UserService::new(user_repo));
        let auth_service = Arc::new(AuthService::new(auth_repo, env.clone()));
        Self {
            env,
            pool,
            user_service,
            auth_service,
        }
    }
}

/// Build the axum router (without binding it).
pub fn build_router(state: AppState) -> axum::Router {
    crate::rest::build(state)
}

/// Boot the HTTP server end-to-end. Blocks until SIGINT/SIGTERM.
pub async fn serve() -> Result<()> {
    let _ = dotenvy::dotenv();
    let env = Env::from_environment().context("load env")?;
    init_logging(&env);

    info!(
        service_name = %env.service_name,
        environment = %env.app_environment,
        "go-clean starting"
    );

    let state = AppState::from_env(env.clone()).await?;
    let router = build_router(state);

    let host: std::net::IpAddr = env
        .app_host
        .parse()
        .unwrap_or(std::net::IpAddr::V4(std::net::Ipv4Addr::LOCALHOST));
    let addr = SocketAddr::new(host, env.app_port);
    let listener = TcpListener::bind(addr)
        .await
        .with_context(|| format!("bind {addr}"))?;
    info!(%addr, "go-clean listening");

    axum::serve(
        listener,
        router.into_make_service_with_connect_info::<SocketAddr>(),
    )
    .with_graceful_shutdown(shutdown_signal())
    .await
    .context("axum serve")?;

    info!("go-clean shut down cleanly");
    Ok(())
}

fn init_logging(env: &Env) {
    use tracing_subscriber::EnvFilter;
    use tracing_subscriber::layer::SubscriberExt;
    use tracing_subscriber::util::SubscriberInitExt;

    let default_level = if env.is_production() { "info" } else { "debug" };
    let filter =
        EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new(default_level));

    let _ = tracing_subscriber::registry()
        .with(filter)
        .with(tracing_subscriber::fmt::layer().with_target(true))
        .try_init();
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c().await.ok();
    };
    #[cfg(unix)]
    let terminate = async {
        if let Ok(mut sig) = signal::unix::signal(signal::unix::SignalKind::terminate()) {
            sig.recv().await;
        }
    };
    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();
    tokio::select! {
        () = ctrl_c => {}
        () = terminate => {}
    }
    info!("shutdown signal received");
}
