//! Shared test harness for go-modular integration tests.
//!
//! Each test file under `tests/` pulls this in via `mod common;`
//! and a `use common::*;` (or named imports). Because `tests/common/mod.rs`
//! is included by each test binary separately, items not used in a
//! given binary would normally trigger `dead_code` warnings —
//! `#![allow(dead_code)]` at the crate root of the common module
//! suppresses those.
//!
//! Provides:
//! - `test_config(db_url)` — builds a test `Config` with the
//!   mailer pinned to noop mode.
//! - `spawn_state()` — starts a disposable Postgres 16-alpine
//!   testcontainer, runs migrations, returns the wired `AppState`
//!   plus the container guard (MUST be held for the test lifetime).
//! - `test_router(state)` — builds the axum router with
//!   `MockConnectInfo` so `ConnectInfo<SocketAddr>` extractors
//!   resolve without a real listener.
//! - `call(router, request)` — one-shot HTTP dispatch that returns
//!   `(StatusCode, serde_json::Value)`.
//! - `json_request` / `authed_request` — request builders.
//! - `insert_verified_user(state, name, email, password)` — seeds
//!   a verified user + argon2id password row.
//! - `signin(router, email, password)` — signs the user in via
//!   the HTTP endpoint and returns the decoded response.
//! - `extract_tokens(body)` — parses a signin / rotation response
//!   into `(access, refresh, session_id)`.

#![allow(dead_code)]

use std::net::SocketAddr;

use axum::Router;
use axum::body::Body;
use axum::extract::connect_info::MockConnectInfo;
use axum::http::{Method, Request, StatusCode};
use http_body_util::BodyExt;
use serde_json::{Value, json};
use sqlx::postgres::PgPoolOptions;
use testcontainers_modules::postgres::Postgres;
use testcontainers_modules::testcontainers::runners::AsyncRunner;
use testcontainers_modules::testcontainers::{ContainerAsync, ImageExt};
use tower::ServiceExt;
use uuid::Uuid;

use rust_modular::AppState;
use rust_modular::apputils::PasswordHasher;
use rust_modular::build_router;
use rust_modular::config::Config;

/// Build a test config matching the Phase D defaults with the DB URL
/// pointed at the testcontainer instance. Clears `smtp_host` so the
/// mailer runs in noop mode — the default config points at
/// `localhost:1025` (mailhog) and without a real relay on that port
/// verification initiate would return 500.
pub fn test_config(db_url: &str) -> Config {
    let mut cfg = Config::from_defaults().expect("default config valid");
    cfg.database.database_url = db_url.to_string();
    cfg.app.jwt_secret_key = "integration-test-secret-not-for-prod".to_string();
    cfg.app.app_mode = "test".to_string();
    cfg.app.app_base_url = "http://test.local".to_string();
    cfg.mailer.smtp_host = String::new();
    cfg
}

/// Spin up a Postgres 16-alpine container, apply migrations, return
/// the wired `AppState` + container guard. Hold the guard for the
/// test's lifetime — dropping it stops the container and kills the
/// pool.
pub async fn spawn_state() -> (AppState, ContainerAsync<Postgres>) {
    let container = Postgres::default()
        .with_db_name("rust_modular_test")
        .with_user("postgres")
        .with_password("postgres")
        .with_tag("16-alpine")
        .start()
        .await
        .expect("start postgres 16-alpine container");
    let host_port = container
        .get_host_port_ipv4(5432)
        .await
        .expect("postgres host port");

    let database_url =
        format!("postgres://postgres:postgres@127.0.0.1:{host_port}/rust_modular_test");
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("connect to testcontainer postgres");

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("run migrations against testcontainer");

    let state = AppState::from_parts(test_config(&database_url), pool);
    (state, container)
}

/// Build a test router with `MockConnectInfo` so the `ConnectInfo`
/// extractor in our handlers resolves without a real listener.
pub fn test_router(state: AppState) -> Router {
    build_router(state).layer(MockConnectInfo(SocketAddr::from(([127, 0, 0, 1], 9999))))
}

/// Send a request through the router, parse response body as JSON.
pub async fn call(router: Router, request: Request<Body>) -> (StatusCode, Value) {
    let response = router.oneshot(request).await.expect("router oneshot");
    let status = response.status();
    let body_bytes = response
        .into_body()
        .collect()
        .await
        .expect("collect body")
        .to_bytes();
    let body_json = if body_bytes.is_empty() {
        Value::Null
    } else {
        serde_json::from_slice(&body_bytes).unwrap_or(Value::Null)
    };
    (status, body_json)
}

/// Build a JSON-body request.
pub fn json_request(method: Method, uri: &str, body: &Value) -> Request<Body> {
    Request::builder()
        .method(method)
        .uri(uri)
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(body).unwrap()))
        .unwrap()
}

/// Build a request with an `Authorization: Bearer` header + optional
/// JSON body.
pub fn authed_request(
    method: Method,
    uri: &str,
    token: &str,
    body: Option<&Value>,
) -> Request<Body> {
    let mut builder = Request::builder()
        .method(method)
        .uri(uri)
        .header("authorization", format!("Bearer {token}"));
    if body.is_some() {
        builder = builder.header("content-type", "application/json");
    }
    let body = body.map_or_else(Body::empty, |b| Body::from(serde_json::to_vec(b).unwrap()));
    builder.body(body).unwrap()
}

/// Insert a verified user + argon2id password into the DB, bypassing
/// the HTTP signup flow. Returns the user's UUID.
pub async fn insert_verified_user(
    state: &AppState,
    display_name: &str,
    email: &str,
    password: &str,
) -> Uuid {
    let id = Uuid::now_v7();
    let username = email.split('@').next().unwrap_or("user").replace('.', "_");

    sqlx::query(
        "INSERT INTO public.users \
         (id, display_name, email, username, email_verified_at) \
         VALUES ($1, $2, $3, $4, NOW())",
    )
    .bind(id)
    .bind(display_name)
    .bind(email)
    .bind(&username)
    .execute(&state.pool)
    .await
    .expect("insert user row");

    let hasher = PasswordHasher::new();
    let phc = hasher.hash(password).expect("hash password");
    sqlx::query("INSERT INTO public.user_passwords (user_id, password_hash) VALUES ($1, $2)")
        .bind(id)
        .bind(phc.as_bytes())
        .execute(&state.pool)
        .await
        .expect("insert user_password row");

    id
}

/// Signin via the HTTP endpoint, returning the decoded body.
pub async fn signin(router: Router, email: &str, password: &str) -> (StatusCode, Value) {
    let req = json_request(
        Method::POST,
        "/api/v1/auth/signin/email",
        &json!({ "email": email, "password": password }),
    );
    call(router, req).await
}

/// Extract `(access_token, refresh_token, session_id)` from a
/// successful signin / rotation response.
pub fn extract_tokens(body: &Value) -> (String, String, Uuid) {
    let access = body["access_token"]
        .as_str()
        .expect("access_token present")
        .to_string();
    let refresh = body["refresh_token"]
        .as_str()
        .expect("refresh_token present")
        .to_string();
    let sid = body["session_id"].as_str().expect("session_id present");
    (
        access,
        refresh,
        Uuid::parse_str(sid).expect("session_id uuid"),
    )
}
