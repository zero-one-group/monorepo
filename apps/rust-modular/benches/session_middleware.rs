//! Session-check middleware perf bench (D-IT-4).
//!
//! Measures the end-to-end latency of a protected request going
//! through `require_auth` + session-row DB lookup + the user
//! handler's own query. This is the full user-visible latency of
//! any authenticated endpoint.
//!
//! **Target**: p99 < 5 ms. If the bench misses this, add a
//! `moka::future::Cache` in front of `validate_session()` with a
//! 30-second TTL and re-bench. See migrate-rust.md §10.6 for the
//! documented before/after if the cache was needed.
//!
//! Harness note: criterion benches run as a separate cargo target
//! from `tests/`, so this file can't pull in `tests/common/mod.rs`.
//! The ~40 lines of duplicated setup below are the cost of that
//! separation.

use std::net::SocketAddr;
use std::time::Duration;

use axum::Router;
use axum::body::Body;
use axum::extract::connect_info::MockConnectInfo;
use axum::http::{Method, Request};
use criterion::{Criterion, criterion_group, criterion_main};
use sqlx::postgres::PgPoolOptions;
use testcontainers_modules::postgres::Postgres;
use testcontainers_modules::testcontainers::runners::AsyncRunner;
use testcontainers_modules::testcontainers::{ContainerAsync, ImageExt};
use tokio::runtime::Runtime;
use tower::ServiceExt;

use rust_modular::AppState;
use rust_modular::apputils::PasswordHasher;
use rust_modular::build_router;
use rust_modular::config::Config;

struct BenchCtx {
    router: Router,
    access_token: String,
    #[allow(dead_code)]
    container: ContainerAsync<Postgres>,
}

async fn setup_bench() -> BenchCtx {
    let container = Postgres::default()
        .with_db_name("rust_modular_bench")
        .with_user("postgres")
        .with_password("postgres")
        .with_tag("16-alpine")
        .start()
        .await
        .expect("start postgres container");
    let host_port = container
        .get_host_port_ipv4(5432)
        .await
        .expect("postgres host port");
    let database_url =
        format!("postgres://postgres:postgres@127.0.0.1:{host_port}/rust_modular_bench");

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await
        .expect("connect pool");
    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("apply migrations");

    let mut cfg = Config::from_defaults().expect("default config");
    cfg.database.database_url = database_url;
    cfg.app.jwt_secret_key = "bench-test-secret".to_string();
    cfg.app.app_mode = "test".to_string();
    cfg.app.app_base_url = "http://bench.local".to_string();
    cfg.mailer.smtp_host = String::new();

    let state = AppState::from_parts(cfg, pool);

    // Insert a verified user + password row, then sign in via the
    // HTTP endpoint to get a real access token bound to a real
    // session row (which is what the middleware DB lookup hits).
    let hasher = PasswordHasher::new();
    let phc = hasher.hash("bench password").expect("hash");
    let user_id = uuid::Uuid::now_v7();
    sqlx::query(
        "INSERT INTO public.users \
         (id, display_name, email, username, email_verified_at) \
         VALUES ($1, 'Bench User', 'bench@example.com', 'bench', NOW())",
    )
    .bind(user_id)
    .execute(&state.pool)
    .await
    .expect("insert user");
    sqlx::query("INSERT INTO public.user_passwords (user_id, password_hash) VALUES ($1, $2)")
        .bind(user_id)
        .bind(phc.as_bytes())
        .execute(&state.pool)
        .await
        .expect("insert password");

    let router =
        build_router(state).layer(MockConnectInfo(SocketAddr::from(([127, 0, 0, 1], 9999))));

    let signin_req = Request::builder()
        .method(Method::POST)
        .uri("/api/v1/auth/signin/email")
        .header("content-type", "application/json")
        .body(Body::from(
            serde_json::to_vec(&serde_json::json!({
                "email": "bench@example.com",
                "password": "bench password",
            }))
            .unwrap(),
        ))
        .unwrap();
    let response = router
        .clone()
        .oneshot(signin_req)
        .await
        .expect("signin oneshot");
    let status = response.status();
    let body_bytes = http_body_util::BodyExt::collect(response.into_body())
        .await
        .expect("collect body")
        .to_bytes();
    assert!(
        status.is_success(),
        "signin failed: status={status} body={}",
        String::from_utf8_lossy(&body_bytes)
    );
    let signin: serde_json::Value = serde_json::from_slice(&body_bytes).expect("parse signin");
    let access_token = signin["access_token"]
        .as_str()
        .expect("access_token present")
        .to_string();

    BenchCtx {
        router,
        access_token,
        container,
    }
}

fn bench_session_middleware(c: &mut Criterion) {
    // Build a dedicated multi-thread runtime so both setup and
    // per-iteration awaits can run on the same pool.
    let rt = Runtime::new().expect("tokio runtime");
    let ctx = rt.block_on(setup_bench());

    let mut group = c.benchmark_group("session_middleware");
    group.measurement_time(Duration::from_secs(10));
    group.sample_size(200);

    // **Protected endpoint**: hits require_auth (session-row DB
    // lookup) + the list_users handler (its own SELECT from
    // public.users). This is the user-visible latency of any
    // authenticated endpoint.
    group.bench_function("get_users_authed", |b| {
        b.to_async(&rt).iter(|| {
            let router = ctx.router.clone();
            let token = ctx.access_token.clone();
            async move {
                let req = Request::builder()
                    .method(Method::GET)
                    .uri("/api/v1/users")
                    .header("authorization", format!("Bearer {token}"))
                    .body(Body::empty())
                    .unwrap();
                let resp = router.oneshot(req).await.expect("oneshot");
                assert!(resp.status().is_success(), "status: {}", resp.status());
            }
        });
    });

    // **Baseline**: unprotected `/healthz` for comparison. Any
    // overhead the session-check middleware adds shows up as the
    // delta between this and `get_users_authed`.
    group.bench_function("healthz_baseline", |b| {
        b.to_async(&rt).iter(|| {
            let router = ctx.router.clone();
            async move {
                let req = Request::builder()
                    .method(Method::GET)
                    .uri("/healthz")
                    .body(Body::empty())
                    .unwrap();
                let resp = router.oneshot(req).await.expect("oneshot");
                assert!(resp.status().is_success());
            }
        });
    });

    group.finish();

    // Explicitly drop the ctx (which owns the ContainerAsync)
    // inside the runtime so testcontainers' async-drop impl has
    // a reactor to use. Dropping on the main thread after criterion
    // finishes would panic with "no reactor running".
    rt.block_on(async move { drop(ctx) });
}

criterion_group!(benches, bench_session_middleware);
criterion_main!(benches);
