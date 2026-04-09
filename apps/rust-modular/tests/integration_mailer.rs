//! Mailhog integration test (D-SMTP-4).
//!
//! Spins up a `mailhog/mailhog` testcontainer, configures the
//! `Mailer` to relay on its SMTP port (1025, plaintext for dev),
//! sends a verification email via
//! `mailer.send_verification_email(...)`, then hits mailhog's
//! HTTP API at port 8025 and asserts the email arrived with the
//! expected subject, recipient, and that the body contains the
//! verification URL.
//!
//! Why this test exists: unit tests in `src/mailer/mod.rs` cover
//! the askama template rendering and the noop-mailer Debug impl,
//! but nothing exercised the full lettre transport + SMTP
//! handshake + server-side message inspection. This test closes
//! that gap with a real SMTP relay in a disposable container.

#![allow(clippy::too_many_lines)]

use std::time::Duration;

use reqwest::Client;
use serde_json::Value;
use testcontainers_modules::testcontainers::GenericImage;
use testcontainers_modules::testcontainers::core::{IntoContainerPort, WaitFor};
use testcontainers_modules::testcontainers::runners::AsyncRunner;
use tokio::time::sleep;

use rust_modular::config::MailerConfig;
use rust_modular::mailer::Mailer;

/// Build a Mailhog image definition pointing at the standard
/// SMTP (1025) and HTTP (8025) ports. Mailhog uses Go's `log`
/// package which writes to stderr, so the wait strategy matches
/// `Serving under` on stderr — that's mailhog's final boot line.
fn mailhog_image() -> GenericImage {
    GenericImage::new("mailhog/mailhog", "latest")
        .with_exposed_port(1025.tcp())
        .with_exposed_port(8025.tcp())
        .with_wait_for(WaitFor::message_on_stderr("Serving under"))
}

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn mailer_sends_verification_email_captured_by_mailhog() {
    // 1. Start Mailhog.
    let container = mailhog_image()
        .start()
        .await
        .expect("start mailhog container");
    let smtp_port = container
        .get_host_port_ipv4(1025)
        .await
        .expect("mailhog smtp port");
    let http_port = container
        .get_host_port_ipv4(8025)
        .await
        .expect("mailhog http port");

    // 2. Configure the Mailer to relay through the container.
    //    Port 1025 is plaintext — our Mailer::build_transport
    //    dispatches it to `builder_dangerous` (dev only; the
    //    config validator rejects this in production mode).
    let mailer_cfg = MailerConfig {
        smtp_host: "127.0.0.1".to_string(),
        smtp_port,
        smtp_username: String::new(),
        smtp_password: String::new(),
        smtp_sender_name: "Go-Modular Test".to_string(),
        smtp_sender_email: "noreply@test.local".to_string(),
    };
    let mailer = Mailer::from_config(&mailer_cfg);

    // Mailhog sometimes needs a beat after "Serving under" before
    // it accepts SMTP — give it a short breathing room.
    sleep(Duration::from_millis(250)).await;

    // 3. Send the verification email.
    let to = "alice@test.local";
    let verification_url =
        "https://test.local/api/v1/auth/verify-email?token=mailhog-test-token-42";
    mailer
        .send_verification_email(to, "Alice", verification_url, 15)
        .await
        .expect("send_verification_email");

    // 4. Poll mailhog's HTTP API until the message arrives
    //    (should be near-instant but allow a brief grace period).
    let client = Client::builder()
        .timeout(Duration::from_secs(3))
        .build()
        .expect("reqwest client");
    let api_url = format!("http://127.0.0.1:{http_port}/api/v2/messages");

    let mut messages: Option<Value> = None;
    for _ in 0..20 {
        if let Ok(resp) = client.get(&api_url).send().await
            && let Ok(body) = resp.json::<Value>().await
            && body
                .get("total")
                .and_then(Value::as_u64)
                .is_some_and(|n| n >= 1)
        {
            messages = Some(body);
            break;
        }
        sleep(Duration::from_millis(150)).await;
    }
    let body = messages.expect("mailhog never received the email");

    // 5. Assert the message shape. Mailhog's API v2 response:
    //    { total, count, start, items: [ { From: { Mailbox, Domain },
    //                                      To: [ { Mailbox, Domain } ],
    //                                      Content: { Headers: {...}, Body } } ] }
    assert_eq!(body["total"].as_u64(), Some(1), "body: {body}");

    let item = &body["items"][0];

    // Recipient check.
    let to_mailbox = item["To"][0]["Mailbox"]
        .as_str()
        .expect("To.Mailbox present");
    let to_domain = item["To"][0]["Domain"].as_str().expect("To.Domain present");
    assert_eq!(format!("{to_mailbox}@{to_domain}"), to);

    // Sender check.
    let from_mailbox = item["From"]["Mailbox"]
        .as_str()
        .expect("From.Mailbox present");
    let from_domain = item["From"]["Domain"]
        .as_str()
        .expect("From.Domain present");
    assert_eq!(
        format!("{from_mailbox}@{from_domain}"),
        "noreply@test.local"
    );

    // Subject check — mailhog stores headers as arrays of strings.
    let subject = item["Content"]["Headers"]["Subject"][0]
        .as_str()
        .expect("Subject header present");
    assert!(
        subject.contains("Verify"),
        "subject should mention Verify, got: {subject}"
    );

    // Body check — verification URL must appear in the rendered
    // HTML. Mailhog stores the raw MIME body which uses
    // quoted-printable encoding for long lines, splitting our
    // 65+-char verification URL across multiple lines with
    // `=\r\n` soft-wrap markers. Normalize by removing those
    // markers before substring matching. We do NOT decode the
    // full QP encoding — the substrings we assert on are pure
    // ASCII alphanumerics and don't include any characters that
    // would have been QP-escaped (e.g., `=` → `=3D`).
    let raw_body = item["Content"]["Body"]
        .as_str()
        .expect("Content.Body present");
    let normalized_body = raw_body.replace("=\r\n", "").replace("=\n", "");

    assert!(
        normalized_body.contains("mailhog-test-token-42"),
        "body should contain the verification token; normalized body:\n{normalized_body}"
    );
    assert!(
        normalized_body.contains("Alice"),
        "body should contain the user's display name; normalized body:\n{normalized_body}"
    );
    assert!(
        normalized_body.contains("15 minutes"),
        "body should contain the expiry duration; normalized body:\n{normalized_body}"
    );
}
