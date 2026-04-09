//! Rust port of the Zero One Group {{ package_name | kebab_case }} service.
//!
//! Behavior-equivalent to the original Python `app/` package: same HTTP
//! endpoints, same JSON shapes, same env vars, same `OTel` REPLICABLE
//! spans/metrics. Internals are Rust-idiomatic (axum + tower + sqlx +
//! tokio + async-openai) per the deep-dive spec's Round 1 decision.
//!
//! Public surface for tests and the binary entrypoint:
//! - [`serve`] — boots the HTTP server, blocks until shutdown
//! - [`build_router`] — constructs the [`axum::Router`] alone (used by
//!   integration tests so they can wire up a test client without
//!   binding to a real port)
//! - [`AppState`] — handler-shared state (DB pool, `OpenAI` client, env)

pub mod core;
pub mod model;
pub mod repository;
pub mod router;
pub mod services;

use std::net::SocketAddr;
use std::sync::Arc;

use anyhow::{Context, Result};
use tokio::net::TcpListener;
use tokio::signal;
use tracing::info;

use crate::core::database::Database;
use crate::core::env::Env;
use crate::core::instrumentation::init_observability;
use crate::core::logging::init_logging;
use crate::repository::openai::greeting::GreetingRepoOpenAI;
use crate::services::greeting::GreetingService;

/// Application state shared with every handler via `axum::extract::State`.
#[derive(Clone)]
pub struct AppState {
    pub env: Arc<Env>,
    pub db: Database,
    pub greeting_service: Arc<GreetingService>,
}

impl AppState {
    /// Construct from a loaded [`Env`]. Establishes the database pool and
    /// the `OpenAI` client up-front so handler latency is minimal.
    pub async fn from_env(env: Env) -> Result<Self> {
        let env = Arc::new(env);
        let db = Database::connect(&env.database_url)
            .await
            .context("connect to database")?;
        let mut openai_cfg = async_openai::config::OpenAIConfig::new()
            .with_api_key(env.openai_api_key.clone());
        if let Ok(base) = std::env::var("OPENAI_API_BASE") {
            openai_cfg = openai_cfg.with_api_base(base);
        }
        let openai_client = async_openai::Client::with_config(openai_cfg);
        let repo = Arc::new(GreetingRepoOpenAI::new(openai_client));
        let greeting_service = Arc::new(GreetingService::new(repo));
        Ok(Self {
            env,
            db,
            greeting_service,
        })
    }
}

/// Build the axum router (without binding it).
///
/// Pulled out of [`serve`] so integration tests can construct their own
/// listener and call this exact router with `tower::ServiceExt::oneshot`
/// or `axum_test::TestServer`.
pub fn build_router(state: AppState) -> axum::Router {
    crate::router::build(state)
}

/// Boot the HTTP server end-to-end. Blocks until SIGINT/SIGTERM.
pub async fn serve() -> Result<()> {
    // 1. Load env via figment+dotenvy. Failures here are fatal.
    let _ = dotenvy::dotenv();
    let env = Env::from_environment().context("load env")?;
    init_logging(&env);
    init_observability(&env)?;

    info!(
        app_name = %env.app_name,
        environment = %env.app_environment,
        "{{ package_name | kebab_case }} starting"
    );

    // 2. Build state + router
    let state = AppState::from_env(env.clone()).await?;
    let router = build_router(state.clone());

    // 3. Bind listener
    let port: u16 = std::env::var("PORT")
        .ok()
        .and_then(|s| s.parse().ok())
        .unwrap_or({{ port_number }});
    let host: std::net::IpAddr = std::env::var("HOST")
        .ok()
        .and_then(|s| s.parse().ok())
        .unwrap_or(std::net::IpAddr::V4(std::net::Ipv4Addr::UNSPECIFIED));
    let addr = SocketAddr::new(host, port);
    let listener = TcpListener::bind(addr)
        .await
        .with_context(|| format!("bind {addr}"))?;
    info!(%addr, "{{ package_name | kebab_case }} listening");

    // 4. Serve until shutdown
    axum::serve(listener, router.into_make_service())
        .with_graceful_shutdown(shutdown_signal(state.clone()))
        .await
        .context("axum serve")?;

    info!("{{ package_name | kebab_case }} shut down cleanly");
    Ok(())
}

async fn shutdown_signal(state: AppState) {
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
    info!("shutdown signal received; disposing database connections");
    state.db.dispose().await;
}
