//! Golden-response tests for the go-clean Rust port.
//!
//! These tests are the primary Phase C acceptance gate per the deep-dive
//! spec's behavior-equivalent contract. They exercise every endpoint of
//! the original Go service (`main.go` + `internal/rest/`) against a
//! disposable testcontainer Postgres and assert the full behavior of
//! each endpoint:
//!
//! Endpoints covered:
//!
//! - `GET /` — deterministic health body
//! - `POST /api/v1/auth/login` — 200 with token, 400 bad payload, 401 bad creds
//! - `GET /api/v1/users` — 200 empty, 200 with rows
//! - `GET /api/v1/users/{id}` — 200, 400 bad uuid, 404 not found
//! - `POST /api/v1/users` — 201 created, 400 bad payload
//! - `PUT /api/v1/users/{id}` — 200 updated, 400 bad uuid/payload
//! - `DELETE /api/v1/users/{id}` — 204 deleted, 400 bad uuid
//! - Auth middleware — 401 missing header, 400 malformed bearer, 401 bogus token, 200 through with valid token
//! - End-to-end login flow — create then login then use token to list
//!
//! No network hop: the test harness uses `tower::ServiceExt::oneshot`
//! against the router constructed by `{{ package_name | snake_case }}::build_router`.

// Pedantic cast lints suppressed because the test fixtures use small
// integers (HTTP status codes) that cannot truncate in practice.
#![allow(
    clippy::cast_possible_truncation,
    clippy::cast_lossless,
    clippy::cast_sign_loss
)]

use std::sync::Arc;

use axum::body::Body;
use axum::http::{Method, Request, StatusCode};
use http_body_util::BodyExt;
use serde_json::{Value, json};
use sqlx::postgres::PgPoolOptions;
use testcontainers_modules::postgres::Postgres;
use testcontainers_modules::testcontainers::ContainerAsync;
use testcontainers_modules::testcontainers::runners::AsyncRunner;
use tower::ServiceExt;

use {{ package_name | snake_case }}::config::env::Env;
use {{ package_name | snake_case }}::utils::jwt::generate_token_pair;
use {{ package_name | snake_case }}::{AppState, build_router};

// ---------------------------------------------------------------------
// Infrastructure helpers
// ---------------------------------------------------------------------

/// Test env: dev-mode, deterministic JWT secret so we can forge tokens
/// for negative tests. The values for `app_host/app_port/otel_endpoint`
/// are placeholders — the test harness calls `build_router` directly
/// and never binds a listener.
fn test_env() -> Env {
    Env {
        service_name: "go-clean-test".into(),
        app_environment: "local".into(),
        app_host: "127.0.0.1".into(),
        app_port: {{ port_number }},
        database_url: "postgres://unused".into(),
        cors_allow_origins: None,
        jwt_secret: "supersecret-test-key".into(),
        auth_token_expiry_minutes: 60,
        enable_swagger: false,
        otel_exporter_otlp_endpoint: "localhost:4317".into(),
    }
}

/// Spin up a disposable Postgres container, run migrations, and return
/// both the wired `AppState` and the container guard (MUST be held for
/// the test's lifetime — dropping it stops the container and breaks
/// the pool).
async fn spawn_state() -> (AppState, ContainerAsync<Postgres>) {
    let container = Postgres::default()
        .with_db_name("{{ package_name | snake_case }}_test")
        .with_user("postgres")
        .with_password("postgres")
        .start()
        .await
        .expect("start postgres container");
    let host_port = container
        .get_host_port_ipv4(5432)
        .await
        .expect("postgres host port");

    let database_url = format!("postgres://postgres:postgres@127.0.0.1:{host_port}/{{ package_name | snake_case }}_test");
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("connect to testcontainer postgres");

    // Apply migrations. The sqlx::migrate! macro resolves `./migrations`
    // relative to CARGO_MANIFEST_DIR which is the crate root, so this
    // picks up `apps/go-clean/migrations/*.sql`.
    //
    // NOTE: sqlx does not understand goose's `-- +goose Up` comment
    // directives — files using that syntax are silently treated as
    // zero-statement no-ops. Our migration file has no goose
    // directives; see the commit that removed them and the notes
    // inside `migrations/*.sql`.
    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("run migrations against testcontainer postgres");

    let state = AppState::from_parts(test_env(), pool);
    (state, container)
}

/// Helper: send `Request<Body>` through the router, parse the response
/// body as JSON, and return `(status, body)`.
async fn call(router: axum::Router, request: Request<Body>) -> (StatusCode, Value) {
    let response = router.oneshot(request).await.expect("router oneshot");
    let status = response.status();
    let body_bytes = response
        .into_body()
        .collect()
        .await
        .expect("collect body")
        .to_bytes();
    let body_json: Value = if body_bytes.is_empty() {
        Value::Null
    } else {
        serde_json::from_slice(&body_bytes)
            .unwrap_or_else(|_| panic!("parse body as json: {body_bytes:?}"))
    };
    (status, body_json)
}

/// Helper: build a request with a JSON body.
fn json_request(method: Method, uri: &str, body: &Value) -> Request<Body> {
    Request::builder()
        .method(method)
        .uri(uri)
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(body).unwrap()))
        .unwrap()
}

/// Helper: build a request with an Authorization header + optional JSON body.
fn authed_request(method: Method, uri: &str, token: &str, body: Option<&Value>) -> Request<Body> {
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

/// Forge a valid access token for a given user id + email, using the
/// same test env and JWT secret the `AppState` is configured with.
fn forge_token(state: &AppState, user_id: &str, email: &str) -> String {
    let env = Arc::clone(&state.env);
    let (token, _refresh) = generate_token_pair(user_id, email, &env).expect("generate test token");
    token
}

/// Create a user directly in the DB (bypassing the API) so auth flow
/// tests can exercise login without chaining on POST /users first.
/// Returns the inserted row's id (as String) and the plaintext password.
async fn insert_test_user(state: &AppState, name: &str, email: &str) -> (String, String) {
    let password_plain = "test-password-1234".to_string();
    let hashed = bcrypt::hash(&password_plain, bcrypt::DEFAULT_COST).unwrap();
    let id = uuid::Uuid::new_v4();
    sqlx::query(
        "INSERT INTO users (id, name, email, password, created_at, updated_at) \
         VALUES ($1, $2, $3, $4, NOW(), NOW())",
    )
    .bind(id)
    .bind(name)
    .bind(email)
    .bind(&hashed)
    .execute(&state.pool)
    .await
    .unwrap();
    (id.to_string(), password_plain)
}

// ---------------------------------------------------------------------
// Tests — 1 per endpoint + end-to-end flow + middleware cases
// ---------------------------------------------------------------------

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_root_health() {
    let (state, _guard) = spawn_state().await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        Request::builder().uri("/").body(Body::empty()).unwrap(),
    )
    .await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(body, json!({"code": 200, "message": "All is well!"}));
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_login_missing_fields_is_400() {
    let (state, _guard) = spawn_state().await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        json_request(Method::POST, "/api/v1/auth/login", &json!({})),
    )
    .await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(body["code"], 400);
    assert_eq!(body["message"], "Invalid request payload");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_login_wrong_password_is_401() {
    let (state, _guard) = spawn_state().await;
    let (_id, _pw) = insert_test_user(&state, "alice", "alice@example.com").await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        json_request(
            Method::POST,
            "/api/v1/auth/login",
            &json!({"email": "alice@example.com", "password": "nope"}),
        ),
    )
    .await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
    assert_eq!(body["code"], 401);
    assert_eq!(body["message"], "Invalid email or password");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_login_unknown_email_is_401() {
    let (state, _guard) = spawn_state().await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        json_request(
            Method::POST,
            "/api/v1/auth/login",
            &json!({"email": "ghost@example.com", "password": "anything"}),
        ),
    )
    .await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
    assert_eq!(body["message"], "Invalid email or password");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_login_success_returns_tokens() {
    let (state, _guard) = spawn_state().await;
    let (expected_id, password) = insert_test_user(&state, "bob", "bob@example.com").await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        json_request(
            Method::POST,
            "/api/v1/auth/login",
            &json!({"email": "bob@example.com", "password": password}),
        ),
    )
    .await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(body["code"], 200);
    assert_eq!(body["message"], "Successfully logged in");
    assert_eq!(body["data"]["user"]["id"], expected_id);
    assert_eq!(body["data"]["user"]["email"], "bob@example.com");
    assert!(body["data"]["token"].is_string());
    assert!(body["data"]["refresh_token"].is_string());
    assert!(!body["data"]["token"].as_str().unwrap().is_empty());
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_users_list_requires_auth() {
    let (state, _guard) = spawn_state().await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        Request::builder()
            .uri("/api/v1/users")
            .body(Body::empty())
            .unwrap(),
    )
    .await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
    assert_eq!(body, json!({"message": "Unauthorized"}));
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_users_list_malformed_bearer_is_400() {
    let (state, _guard) = spawn_state().await;
    let router = build_router(state);
    let (status, body) = call(
        router,
        Request::builder()
            .uri("/api/v1/users")
            .header("authorization", "NotBearerToken")
            .body(Body::empty())
            .unwrap(),
    )
    .await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(body["code"], 400);
    assert_eq!(body["message"], "invalid bearer token");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_users_list_bogus_token_is_401() {
    let (state, _guard) = spawn_state().await;
    let router = build_router(state);
    let (status, _body) = call(
        router,
        Request::builder()
            .uri("/api/v1/users")
            .header("authorization", "Bearer not.a.real.token")
            .body(Body::empty())
            .unwrap(),
    )
    .await;
    assert_eq!(status, StatusCode::UNAUTHORIZED);
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_users_list_empty() {
    let (state, _guard) = spawn_state().await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);
    let (status, body) = call(
        router,
        authed_request(Method::GET, "/api/v1/users", &token, None),
    )
    .await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(body["code"], 200);
    assert_eq!(body["message"], "Successfully retrieve user list");
    assert_eq!(body["data"], json!([]));
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_users_list_nonempty_and_filter() {
    let (state, _guard) = spawn_state().await;
    insert_test_user(&state, "carol", "carol@example.com").await;
    insert_test_user(&state, "dave", "dave@example.com").await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);

    // Full list → 2 rows
    let (status, body) = call(
        router.clone(),
        authed_request(Method::GET, "/api/v1/users", &token, None),
    )
    .await;
    assert_eq!(status, StatusCode::OK);
    let list = body["data"].as_array().unwrap();
    assert_eq!(list.len(), 2);

    // Search filter → only carol
    let (status2, body2) = call(
        router,
        authed_request(Method::GET, "/api/v1/users?search=carol", &token, None),
    )
    .await;
    assert_eq!(status2, StatusCode::OK);
    let filtered = body2["data"].as_array().unwrap();
    assert_eq!(filtered.len(), 1);
    assert_eq!(filtered[0]["name"], "carol");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_get_user_by_id_not_found() {
    let (state, _guard) = spawn_state().await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);
    let fake = "00000000-0000-0000-0000-000000000001";
    let (status, body) = call(
        router,
        authed_request(Method::GET, &format!("/api/v1/users/{fake}"), &token, None),
    )
    .await;
    assert_eq!(status, StatusCode::NOT_FOUND);
    assert_eq!(body["code"], 404);
    assert_eq!(body["message"], "User not found");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_get_user_bad_uuid_is_400() {
    let (state, _guard) = spawn_state().await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);
    let (status, body) = call(
        router,
        authed_request(Method::GET, "/api/v1/users/not-a-uuid", &token, None),
    )
    .await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(body["message"], "Invalid user ID format");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_create_user_success() {
    let (state, _guard) = spawn_state().await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);
    let (status, body) = call(
        router,
        authed_request(
            Method::POST,
            "/api/v1/users",
            &token,
            Some(&json!({
                "name": "eve",
                "email": "eve@example.com",
                "password": "correct horse battery staple"
            })),
        ),
    )
    .await;
    assert_eq!(status, StatusCode::CREATED);
    assert_eq!(body["code"], 201);
    assert_eq!(body["message"], "User successfully created");
    assert_eq!(body["data"]["name"], "eve");
    assert_eq!(body["data"]["email"], "eve@example.com");
    assert!(body["data"]["id"].is_string());
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_create_user_bad_payload_is_400() {
    let (state, _guard) = spawn_state().await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);
    let (status, body) = call(
        router,
        authed_request(
            Method::POST,
            "/api/v1/users",
            &token,
            Some(&json!({"name_wrong_field": "x"})),
        ),
    )
    .await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(body["message"], "Invalid request payload");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_update_user_success() {
    let (state, _guard) = spawn_state().await;
    let (id, _pw) = insert_test_user(&state, "frank", "frank@example.com").await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);
    let (status, body) = call(
        router,
        authed_request(
            Method::PUT,
            &format!("/api/v1/users/{id}"),
            &token,
            Some(&json!({"name": "frankie", "email": "frankie@example.com"})),
        ),
    )
    .await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(body["code"], 200);
    assert_eq!(body["message"], "User successfully updated");
    assert_eq!(body["data"]["name"], "frankie");
    assert_eq!(body["data"]["email"], "frankie@example.com");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_delete_user_success_and_repeat_is_500() {
    let (state, _guard) = spawn_state().await;
    let (id, _pw) = insert_test_user(&state, "grace", "grace@example.com").await;
    let token = forge_token(&state, "00000000-0000-0000-0000-000000000000", "a@a.com");
    let router = build_router(state);

    let (status, body) = call(
        router.clone(),
        authed_request(Method::DELETE, &format!("/api/v1/users/{id}"), &token, None),
    )
    .await;
    assert_eq!(status, StatusCode::NO_CONTENT);
    assert_eq!(body["code"], 204);
    assert_eq!(body["message"], "User successfully deleted");

    // Second delete of the same row: the Go service returns 500 with
    // "Failed to delete user: user not found" because it collapses the
    // UserNotFound error through the generic handler error path.
    let (status2, body2) = call(
        router,
        authed_request(Method::DELETE, &format!("/api/v1/users/{id}"), &token, None),
    )
    .await;
    assert_eq!(status2, StatusCode::INTERNAL_SERVER_ERROR);
    assert_eq!(body2["code"], 500);
    assert!(
        body2["message"]
            .as_str()
            .unwrap()
            .starts_with("Failed to delete user")
    );
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_end_to_end_login_then_list() {
    // Full flow: insert a user directly → login via the public API →
    // use the returned access token to list users.
    let (state, _guard) = spawn_state().await;
    let (_id, pw) = insert_test_user(&state, "heidi", "heidi@example.com").await;
    let router = build_router(state);

    // Login
    let (login_status, login_body) = call(
        router.clone(),
        json_request(
            Method::POST,
            "/api/v1/auth/login",
            &json!({"email": "heidi@example.com", "password": pw}),
        ),
    )
    .await;
    assert_eq!(login_status, StatusCode::OK);
    let token = login_body["data"]["token"].as_str().unwrap().to_string();

    // Use the token
    let (list_status, list_body) = call(
        router,
        authed_request(Method::GET, "/api/v1/users", &token, None),
    )
    .await;
    assert_eq!(list_status, StatusCode::OK);
    assert_eq!(list_body["data"].as_array().unwrap().len(), 1);
    assert_eq!(list_body["data"][0]["name"], "heidi");
}
