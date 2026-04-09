//! Live-boot smoke test (D-IT-5).
//!
//! Spawns the real `{{ package_name | kebab_case }}` binary as a subprocess pointing at a
//! fresh Postgres 16-alpine testcontainer, runs migrations via
//! `{{ package_name | kebab_case }} migrate run`, starts `{{ package_name | kebab_case }} serve` on a random
//! port, then hits `/healthz` and `/api/openapi.json` with a real
//! HTTP client. Asserts both endpoints respond and the `OpenAPI` body
//! is valid JSON containing the expected metadata.
//!
//! Unlike `integration_auth.rs` (which uses
//! `tower::ServiceExt::oneshot` against the in-process router),
//! this test exercises the full binary including clap CLI dispatch,
//! dotenvy, figment config loading, the real tokio runtime, and
//! the axum HTTP listener. That's the point — it catches issues
//! like missing env vars, broken `serve()` plumbing, or panics
//! during boot that the in-process tests can't see.
//!
//! Cleanup is deterministic: the server subprocess is `.kill()`-ed
//! on drop via the `ServerGuard` RAII wrapper, and the Postgres
//! container is dropped at end of scope by testcontainers.

#![allow(clippy::too_many_lines)]

use std::net::TcpListener;
use std::process::Stdio;
use std::time::Duration;

use reqwest::Client;
use serde_json::Value;
use testcontainers_modules::postgres::Postgres;
use testcontainers_modules::testcontainers::runners::AsyncRunner;
use testcontainers_modules::testcontainers::{ContainerAsync, ImageExt};
use tokio::process::{Child, Command};
use tokio::time::sleep;

/// RAII wrapper that kills the child process when dropped.
struct ServerGuard {
    child: Option<Child>,
}

impl ServerGuard {
    fn new(child: Child) -> Self {
        Self { child: Some(child) }
    }
}

impl Drop for ServerGuard {
    fn drop(&mut self) {
        if let Some(mut child) = self.child.take() {
            // Best-effort kill; ignore errors if the process already exited.
            let _ = child.start_kill();
        }
    }
}

/// Pick a random available TCP port on 127.0.0.1.
fn pick_free_port() -> u16 {
    let listener = TcpListener::bind("127.0.0.1:0").expect("bind ephemeral port");
    let port = listener.local_addr().expect("local_addr").port();
    drop(listener);
    port
}

/// Common env vars for both `migrate run` and `serve`.
fn base_env(database_url: &str, port: u16) -> Vec<(&'static str, String)> {
    vec![
        ("DATABASE_URL", database_url.to_string()),
        ("APP_MODE", "test".to_string()),
        ("APP_BASE_URL", format!("http://127.0.0.1:{port}")),
        (
            "JWT_SECRET_KEY",
            "smoke-test-secret-not-for-prod".to_string(),
        ),
        ("SERVER_HOST", "127.0.0.1".to_string()),
        ("SERVER_PORT", port.to_string()),
        // Noop mailer — no real SMTP relay in the test environment.
        ("SMTP_HOST", String::new()),
    ]
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn smoke_boot_serves_healthz_and_openapi_json() {
    // 1. Start the Postgres container (matches integration_auth.rs).
    let container: ContainerAsync<Postgres> = Postgres::default()
        .with_db_name("{{ package_name | snake_case }}_smoke")
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
        format!("postgres://postgres:postgres@127.0.0.1:{host_port}/{{ package_name | snake_case }}_smoke");

    // 2. Pick a random port for the server to bind.
    let server_port = pick_free_port();
    let env = base_env(&database_url, server_port);

    // 3. Run migrations via `{{ package_name | kebab_case }} migrate run`. Uses the
    //    `CARGO_BIN_EXE_{{ package_name | kebab_case }}` env var that cargo sets for
    //    integration tests — points at the freshly-built binary.
    let bin = env!("CARGO_BIN_EXE_{{ package_name | kebab_case }}");
    let migrate_output = Command::new(bin)
        .envs(env.iter().map(|(k, v)| (*k, v.as_str())))
        .arg("migrate")
        .arg("run")
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .output()
        .await
        .expect("spawn migrate run");
    assert!(
        migrate_output.status.success(),
        "migrate run failed\nstdout: {}\nstderr: {}",
        String::from_utf8_lossy(&migrate_output.stdout),
        String::from_utf8_lossy(&migrate_output.stderr)
    );

    // 4. Start the server as `{{ package_name | kebab_case }} serve` subprocess.
    let child = Command::new(bin)
        .envs(env.iter().map(|(k, v)| (*k, v.as_str())))
        .arg("serve")
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .kill_on_drop(true)
        .spawn()
        .expect("spawn {{ package_name | kebab_case }} serve");
    let _guard = ServerGuard::new(child);

    // 5. Wait for the server to become ready (up to ~10 seconds).
    let client = Client::builder()
        .timeout(Duration::from_secs(3))
        .build()
        .expect("build reqwest client");
    let base_url = format!("http://127.0.0.1:{server_port}");

    let mut ready = false;
    for _ in 0..30 {
        if let Ok(resp) = client.get(format!("{base_url}/healthz")).send().await
            && resp.status().is_success()
        {
            ready = true;
            break;
        }
        sleep(Duration::from_millis(400)).await;
    }
    assert!(ready, "server never became healthy on {base_url}/healthz");

    // 6. Assert /healthz returns {"status": "up"} with a DB check.
    let resp = client
        .get(format!("{base_url}/healthz"))
        .send()
        .await
        .expect("GET /healthz");
    assert_eq!(resp.status(), 200);
    let body: Value = resp.json().await.expect("healthz JSON");
    assert_eq!(body["status"], "up", "healthz body: {body}");
    assert_eq!(body["details"]["database"]["status"], "up");

    // 7. Assert /api/openapi.json is valid JSON with the expected
    //    metadata and at least 14 path entries (utoipa collapses
    //    multiple HTTP methods onto a single path).
    let resp = client
        .get(format!("{base_url}/api/openapi.json"))
        .send()
        .await
        .expect("GET /api/openapi.json");
    assert_eq!(resp.status(), 200);
    let body: Value = resp.json().await.expect("openapi JSON");
    assert_eq!(body["info"]["title"], "{{ package_name | kebab_case }}");
    assert_eq!(body["openapi"], "3.1.0");

    let paths = body["paths"]
        .as_object()
        .expect("openapi paths object present");
    assert!(
        paths.len() >= 14,
        "expected at least 14 unique paths, got {}: {:?}",
        paths.len(),
        paths.keys().collect::<Vec<_>>()
    );
    // Spot-check a representative route.
    assert!(paths.contains_key("/api/v1/auth/signin/email"));
    assert!(paths.contains_key("/api/v1/users/{userId}"));

    // 8. ServerGuard drops → SIGKILL the subprocess. Postgres
    //    container drops at scope exit.
}
