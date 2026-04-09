//! Golden-response tests for the fastapi-ai Rust port.
//!
//! These tests are the primary acceptance gate for Phase B. They load
//! `tests/golden/fixtures.json` (captured from the live Python service
//! on 2026-04-09) and replay every fixture against the Rust router,
//! asserting the result matches per-fixture `match_mode`:
//!
//! - `"exact"` — response body deserialized as a `serde_json::Value`
//!   must equal the fixture's `expected_body` key-for-key. Used for the
//!   deterministic `/` and `/health-check` endpoints.
//! - `"schema"` — response body must satisfy the fixture's
//!   `schema_assertions` block. Used for `/openai/greetings`, whose
//!   LLM-produced content varies per call but whose envelope
//!   (`success`, `message`, `data.response.greetings[*].{language, greeting}`)
//!   is stable.
//!
//! Infrastructure:
//! - Postgres: `testcontainers-modules` spins up a disposable
//!   `postgres:16-alpine` container. Migrations are applied with
//!   `sqlx::migrate!("./migrations")`.
//! - `OpenAI`: `wiremock` stands up a local HTTP mock server pretending
//!   to be the `OpenAI` Chat Completions API. The async-openai client
//!   is pointed at it via `OpenAIConfig::with_api_base(...)`. The
//!   mock returns a canned Chat Completions response whose message
//!   content is a valid JSON object matching the Python service's
//!   `{"greetings": [...]}` contract.
//! - Router: called via `tower::ServiceExt::oneshot` — no listener
//!   bound, no network hop — the fastest possible path that still
//!   exercises the full handler + middleware stack.

// Pedantic cast lints triggered by narrowing u64 (serde_json numeric
// values) to u16 (HTTP status) or usize (array index bounds). Every
// affected value comes from `fixtures.json`, which is hand-authored
// and contains small integers (status codes, assertion bounds ≤ 32),
// so the casts cannot truncate in practice.
#![allow(
    clippy::cast_possible_truncation,
    clippy::cast_lossless,
    clippy::cast_sign_loss
)]

use std::sync::Arc;

use async_openai::Client;
use async_openai::config::OpenAIConfig;
use axum::body::Body;
use axum::http::{Request, StatusCode};
use http_body_util::BodyExt;
use serde_json::{Value, json};
use sqlx::postgres::PgPoolOptions;
use testcontainers_modules::postgres::Postgres;
use testcontainers_modules::testcontainers::runners::AsyncRunner;
use tower::ServiceExt;
use wiremock::matchers::{method, path};
use wiremock::{Mock, MockServer, ResponseTemplate};

use rust_ai::core::database::Database;
use rust_ai::core::env::Env;
use rust_ai::core::instrumentation::init_observability;
use rust_ai::repository::openai::greeting::GreetingRepoOpenAI;
use rust_ai::services::greeting::GreetingService;
use rust_ai::{AppState, build_router};

// ---------------------------------------------------------------------
// Fixture loader
// ---------------------------------------------------------------------

/// Fixture file path, relative to the crate root.
const FIXTURES_PATH: &str = "tests/golden/fixtures.json";

fn load_fixtures() -> Value {
    let raw = std::fs::read_to_string(FIXTURES_PATH).expect("read fixtures.json");
    serde_json::from_str(&raw).expect("parse fixtures.json")
}

fn fixture<'a>(fixtures: &'a Value, name: &str) -> &'a Value {
    fixtures["fixtures"]
        .as_array()
        .expect("fixtures array")
        .iter()
        .find(|f| f["name"].as_str() == Some(name))
        .unwrap_or_else(|| panic!("fixture `{name}` not found in {FIXTURES_PATH}"))
}

// ---------------------------------------------------------------------
// Infrastructure helpers
// ---------------------------------------------------------------------

/// Spin up a disposable Postgres container, run migrations, and return
/// the wired `Database`. The returned guard MUST be held for the test's
/// lifetime — dropping it stops the container and breaks the pool.
async fn spawn_postgres() -> (
    Database,
    testcontainers_modules::testcontainers::ContainerAsync<Postgres>,
) {
    let container = Postgres::default()
        .start()
        .await
        .expect("start postgres container");
    let host_port = container
        .get_host_port_ipv4(5432)
        .await
        .expect("postgres host port");

    let database_url = format!("postgres://postgres:postgres@127.0.0.1:{host_port}/postgres");
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("connect to testcontainer postgres");

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("run migrations against testcontainer postgres");

    (Database::from_pool(pool), container)
}

/// Spin up a wiremock mock server pretending to be the `OpenAI` Chat
/// Completions endpoint. The canned response is a valid Chat Completions
/// payload whose `content` field is a JSON-object string matching the
/// Python service's `{"greetings": [{"language":..., "greeting":...}]}` shape.
async fn spawn_openai_mock() -> MockServer {
    let mock = MockServer::start().await;

    let canned_chat_completions_body = json!({
        "id": "chatcmpl-test-0001",
        "object": "chat.completion",
        "created": 1_712_000_000,
        "model": "gpt-4.1-mini",
        "choices": [
            {
                "index": 0,
                "message": {
                    "role": "assistant",
                    "content": "{\"greetings\": [{\"language\": \"English\", \"greeting\": \"Hello\"}, {\"language\": \"Spanish\", \"greeting\": \"Hola\"}, {\"language\": \"French\", \"greeting\": \"Bonjour\"}, {\"language\": \"German\", \"greeting\": \"Hallo\"}, {\"language\": \"Japanese\", \"greeting\": \"\\u3053\\u3093\\u306b\\u3061\\u306f\"}]}",
                    "refusal": null
                },
                "logprobs": null,
                "finish_reason": "stop"
            }
        ],
        "usage": {
            "prompt_tokens": 42,
            "completion_tokens": 84,
            "total_tokens": 126
        }
    });

    Mock::given(method("POST"))
        .and(path("/chat/completions"))
        .respond_with(ResponseTemplate::new(200).set_body_json(canned_chat_completions_body))
        .mount(&mock)
        .await;

    mock
}

/// Build a minimal test `Env` with fake secrets. Values are placeholders
/// — the real database URL + `OpenAI` base URL are injected through the
/// `AppState` fields, not through this env.
fn test_env() -> Env {
    Env {
        ml_prefix_api: String::new(),
        app_name: "fastapi-ai-test".into(),
        app_environment: "development".into(),
        database_url: "postgres://unused".into(),
        openai_api_key: "test-key".into(),
        otel_exporter_otlp_endpoint: "localhost:4317".into(),
    }
}

/// Wire an `AppState` from a live testcontainer `Database` and a
/// wiremock-backed `OpenAI` client.
fn build_test_state(db: Database, openai_base: &str) -> AppState {
    // Ensure the metrics registry is initialized. `init_observability`
    // is idempotent via `OnceLock::set`.
    init_observability(&test_env()).expect("init observability");

    let client = Client::with_config(
        OpenAIConfig::new()
            .with_api_key("test-key")
            .with_api_base(openai_base.to_string()),
    );
    let repo = Arc::new(GreetingRepoOpenAI::new(client));
    let greeting_service = Arc::new(GreetingService::new(repo));

    AppState {
        env: Arc::new(test_env()),
        db,
        greeting_service,
    }
}

/// Send a `Request<Body>` through the router and return the parsed JSON
/// body plus the status code.
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
        serde_json::from_slice(&body_bytes).expect("parse response body as json")
    };
    (status, body_json)
}

// ---------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_root_exact_match() {
    let fixtures = load_fixtures();
    let fx = fixture(&fixtures, "root");
    assert_eq!(fx["match_mode"], "exact");

    let (db, _container) = spawn_postgres().await;
    let mock = spawn_openai_mock().await;
    let state = build_test_state(db, &mock.uri());
    let router = build_router(state);

    let (status, body) = call(
        router,
        Request::builder().uri("/").body(Body::empty()).unwrap(),
    )
    .await;

    assert_eq!(
        status.as_u16(),
        fx["expected_status"].as_u64().unwrap() as u16
    );
    assert_eq!(body, fx["expected_body"], "body mismatch for `root`");
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_health_check_ok_exact_match() {
    let fixtures = load_fixtures();
    let fx = fixture(&fixtures, "health_check_ok");
    assert_eq!(fx["match_mode"], "exact");

    let (db, _container) = spawn_postgres().await;
    let mock = spawn_openai_mock().await;
    let state = build_test_state(db, &mock.uri());
    let router = build_router(state);

    let (status, body) = call(
        router,
        Request::builder()
            .uri("/health-check")
            .body(Body::empty())
            .unwrap(),
    )
    .await;

    assert_eq!(
        status.as_u16(),
        fx["expected_status"].as_u64().unwrap() as u16
    );
    assert_eq!(
        body, fx["expected_body"],
        "body mismatch for `health_check_ok`"
    );
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_openai_greetings_schema_match() {
    let fixtures = load_fixtures();
    let fx = fixture(&fixtures, "openai_greetings");
    assert_eq!(fx["match_mode"], "schema");

    let (db, _container) = spawn_postgres().await;
    let mock = spawn_openai_mock().await;
    let state = build_test_state(db, &mock.uri());
    let router = build_router(state);

    let (status, body) = call(
        router,
        Request::builder()
            .uri("/openai/greetings")
            .body(Body::empty())
            .unwrap(),
    )
    .await;

    assert_eq!(
        status.as_u16(),
        fx["expected_status"].as_u64().unwrap() as u16
    );

    let assertions = &fx["schema_assertions"];

    // Top-level keys
    let expected_top_keys: Vec<&str> = assertions["top_level_keys"]
        .as_array()
        .expect("top_level_keys array")
        .iter()
        .map(|v| v.as_str().expect("top level key is string"))
        .collect();
    for key in &expected_top_keys {
        assert!(
            body.get(*key).is_some(),
            "response body missing top-level key `{key}`; body = {body}"
        );
    }

    // success == true
    assert_eq!(
        body["success"], assertions["success_is_true"],
        "success flag mismatch"
    );

    // message == "Operation successful"
    assert_eq!(
        body["message"], assertions["message_equals"],
        "message mismatch"
    );

    // data.response.greetings is an array
    let greetings = body["data"]["response"]["greetings"]
        .as_array()
        .expect("data.response.greetings must be an array");

    // bounds check
    let min = assertions["greetings_min_length"].as_u64().unwrap() as usize;
    let max = assertions["greetings_max_length"].as_u64().unwrap() as usize;
    assert!(
        greetings.len() >= min && greetings.len() <= max,
        "greetings length {} not in [{}, {}]",
        greetings.len(),
        min,
        max
    );

    // Each item has the required keys AND both values are strings
    let required_keys: Vec<&str> = assertions["greeting_item_required_keys"]
        .as_array()
        .expect("greeting_item_required_keys array")
        .iter()
        .map(|v| v.as_str().expect("required key is string"))
        .collect();
    for (i, item) in greetings.iter().enumerate() {
        for key in &required_keys {
            let field = item.get(*key).unwrap_or_else(|| {
                panic!("greeting item {i} missing required key `{key}`; item = {item}")
            });
            assert!(
                field.is_string(),
                "greeting item {i} key `{key}` is not a string: {field}"
            );
            assert!(
                !field.as_str().unwrap().is_empty(),
                "greeting item {i} key `{key}` is an empty string"
            );
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn golden_health_check_failure_returns_503() {
    // Complement to golden_health_check_ok_exact_match: when the DB is
    // unreachable, the Python service returns `503 {"detail": "Database
    // connection error"}`. The Rust port must match that behavior.
    //
    // We simulate an unreachable DB by creating a pool pointing at a
    // port nothing is listening on, then closing it.
    let unreachable_pool = PgPoolOptions::new()
        .max_connections(1)
        .acquire_timeout(std::time::Duration::from_millis(200))
        .connect_lazy("postgres://postgres:postgres@127.0.0.1:1/doesnotexist")
        .expect("build lazy pool");
    let db = Database::from_pool(unreachable_pool);

    let mock = spawn_openai_mock().await;
    let state = build_test_state(db, &mock.uri());
    let router = build_router(state);

    let (status, body) = call(
        router,
        Request::builder()
            .uri("/health-check")
            .body(Body::empty())
            .unwrap(),
    )
    .await;

    assert_eq!(status, StatusCode::SERVICE_UNAVAILABLE);
    assert_eq!(body, json!({"detail": "Database connection error"}));
}
