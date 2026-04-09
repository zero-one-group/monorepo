//! HTTP server module — builds the router, binds the listener,
//! serves with graceful shutdown.
//!
//! Ports `internal/server/{server,loader,handler}.go`. The loader's
//! role is absorbed into [`AppState::from_config`]; the handler file
//! becomes [`handler::healthz`] / [`handler::api_docs`].
//!
//! At D-INFRA-11 the router only mounts infrastructure routes
//! (healthz, api-docs, openapi.json) + a 404 catch-all. Module routes
//! from D-USER-* and D-AUTH-* land in [`router::build_router`] as
//! those phases complete.

pub mod handler;
pub mod router;

use std::net::SocketAddr;

use anyhow::{Context, Result};
use tokio::net::TcpListener;
use tokio::signal;
use tracing::info;

use crate::AppState;
use crate::config::Config;

/// Boot the HTTP server end-to-end. Blocks until SIGINT/SIGTERM.
pub async fn serve() -> Result<()> {
    let _ = dotenvy::dotenv();
    let config = Config::from_environment().context("load config")?;
    crate::observer::init_tracing(&config.logging, &config.otel);

    info!(
        service_name = %config.otel.otel_service_name,
        app_mode = %config.app.app_mode,
        "go-modular starting"
    );

    let state = AppState::from_config(config.clone()).await?;
    let app = router::build_router(state);

    let host: std::net::IpAddr = config
        .app
        .server_host
        .parse()
        .unwrap_or(std::net::IpAddr::V4(std::net::Ipv4Addr::UNSPECIFIED));
    let addr = SocketAddr::new(host, config.app.server_port);
    let listener = TcpListener::bind(addr)
        .await
        .with_context(|| format!("bind {addr}"))?;
    info!(%addr, "go-modular listening");

    axum::serve(
        listener,
        app.into_make_service_with_connect_info::<SocketAddr>(),
    )
    .with_graceful_shutdown(shutdown_signal())
    .await
    .context("axum serve")?;

    crate::observer::shutdown_tracing();
    info!("go-modular shut down cleanly");
    Ok(())
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
