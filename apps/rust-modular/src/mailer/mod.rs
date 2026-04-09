//! lettre SMTP mailer + askama email templates (D-SMTP-1..4).
//!
//! **Port-based TLS selection** (design 3.5, replaces the deleted
//! `SMTPSecure` bool config):
//!
//! | Port  | Transport               | Allowed in production? |
//! |-------|-------------------------|------------------------|
//! | 587   | STARTTLS (required)     | Yes                    |
//! | 465   | Implicit TLS (wrapped)  | Yes                    |
//! | 25    | Plaintext               | No (validator rejects) |
//! | 1025  | Plaintext (mailhog dev) | No (validator rejects) |
//!
//! Inline askama templates keep the Rust crate's template source
//! separate from Go's `apps/go-modular/templates/` tree (which gets
//! `git rm`-ed in D-REL-1).
//!
//! **Graceful fallback**: if the transport can't be built (no host,
//! init error), `Mailer::noop()` is used and `send_email` logs the
//! payload at WARN so dev flows still work without a real relay.
//! Matches the Go `svc_verification.go:374` behavior.

use std::sync::Arc;

use anyhow::Result;
use askama::Template;
use lettre::message::Mailbox;
use lettre::transport::smtp::AsyncSmtpTransport;
use lettre::transport::smtp::authentication::Credentials;
use lettre::{AsyncTransport, Message, Tokio1Executor};
use tracing::{info, warn};

use crate::config::MailerConfig;
use crate::domain::AppError;

/// SMTP mailer with optional transport.
///
/// `transport = None` means the mailer is in noop mode: every
/// `send_email` call logs the payload at WARN instead of relaying.
#[derive(Clone)]
pub struct Mailer {
    transport: Option<Arc<AsyncSmtpTransport<Tokio1Executor>>>,
    sender: Mailbox,
}

impl std::fmt::Debug for Mailer {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("Mailer")
            .field("sender", &self.sender)
            .field(
                "transport",
                &if self.transport.is_some() {
                    "<configured>"
                } else {
                    "noop"
                },
            )
            .finish()
    }
}

impl Mailer {
    /// Build a mailer from the loaded config.
    ///
    /// Returns a noop mailer (with a WARN log) instead of an error
    /// when transport init fails — dev flows must still work without
    /// a real SMTP relay available.
    #[must_use]
    pub fn from_config(cfg: &MailerConfig) -> Self {
        let sender_addr = cfg
            .smtp_sender_email
            .parse::<lettre::Address>()
            .ok()
            .unwrap_or_else(|| {
                warn!(
                    "mailer: invalid SMTP_SENDER_EMAIL `{}`; using noreply@localhost",
                    cfg.smtp_sender_email
                );
                "noreply@localhost".parse().expect("static fallback")
            });
        let sender = Mailbox::new(
            Some(cfg.smtp_sender_name.clone()).filter(|s| !s.is_empty()),
            sender_addr,
        );

        let transport = Self::build_transport(cfg);
        if transport.is_none() {
            warn!(
                smtp_host = %cfg.smtp_host,
                smtp_port = cfg.smtp_port,
                "mailer: no transport configured (noop mode; emails logged instead of sent)"
            );
        }

        Self {
            transport: transport.map(Arc::new),
            sender,
        }
    }

    /// Explicit noop constructor (used by unit tests).
    #[must_use]
    pub fn noop() -> Self {
        Self {
            transport: None,
            sender: "noreply@localhost"
                .parse::<Mailbox>()
                .expect("static mailbox"),
        }
    }

    fn build_transport(cfg: &MailerConfig) -> Option<AsyncSmtpTransport<Tokio1Executor>> {
        let host = cfg.smtp_host.trim();
        if host.is_empty() {
            return None;
        }

        // Build credentials if both fields are set.
        let creds = if !cfg.smtp_username.is_empty() && !cfg.smtp_password.is_empty() {
            Some(Credentials::new(
                cfg.smtp_username.clone(),
                cfg.smtp_password.clone(),
            ))
        } else {
            None
        };

        // Localhost plaintext escape hatch: when the host is a
        // loopback address AND the port is NOT a standard TLS
        // port (587/465), fall through to plaintext. This
        // supports dev tooling like mailhog that listens on
        // random host-mapped ports under testcontainers.
        //
        // Production safety: the config validator rejects
        // plaintext SMTP ports (25/1025) in production mode, and
        // a local operator running on 587 or 465 still gets
        // proper TLS via the regular match arms below.
        let is_loopback = matches!(host, "127.0.0.1" | "localhost" | "::1");
        let is_tls_port = matches!(cfg.smtp_port, 587 | 465);

        let builder_result = if is_loopback && !is_tls_port {
            Ok(AsyncSmtpTransport::<Tokio1Executor>::builder_dangerous(host).port(cfg.smtp_port))
        } else {
            match cfg.smtp_port {
                587 => {
                    AsyncSmtpTransport::<Tokio1Executor>::starttls_relay(host).map(|b| b.port(587))
                }
                465 => AsyncSmtpTransport::<Tokio1Executor>::relay(host).map(|b| b.port(465)),
                25 | 1025 => Ok(
                    AsyncSmtpTransport::<Tokio1Executor>::builder_dangerous(host)
                        .port(cfg.smtp_port),
                ),
                other => {
                    // Default to STARTTLS for any other port — safest choice.
                    AsyncSmtpTransport::<Tokio1Executor>::starttls_relay(host)
                        .map(|b| b.port(other))
                }
            }
        };

        let mut builder = builder_result.ok()?;
        if let Some(c) = creds {
            builder = builder.credentials(c);
        }
        Some(builder.build())
    }

    /// Send an email with a rendered HTML body.
    pub async fn send_email(
        &self,
        to: &str,
        subject: &str,
        html_body: String,
    ) -> Result<(), AppError> {
        // Noop mode: log and return success so the caller's flow
        // (verification initiate, etc.) still proceeds in dev.
        let Some(transport) = self.transport.as_ref() else {
            info!(to, subject, "mailer: noop — would have sent email");
            return Ok(());
        };

        let to_addr: Mailbox = to
            .parse()
            .map_err(|e| AppError::BadRequest(format!("invalid recipient: {e}")))?;

        let message = Message::builder()
            .from(self.sender.clone())
            .to(to_addr)
            .subject(subject)
            .header(lettre::message::header::ContentType::TEXT_HTML)
            .body(html_body)
            .map_err(|e| AppError::Internal(anyhow::anyhow!("build message: {e}")))?;

        transport
            .send(message)
            .await
            .map_err(|e| AppError::Internal(anyhow::anyhow!("smtp send: {e}")))?;
        Ok(())
    }

    /// Render the verification email template and send it.
    pub async fn send_verification_email(
        &self,
        to: &str,
        user_name: &str,
        verification_url: &str,
        expiry_minutes: u64,
    ) -> Result<(), AppError> {
        let tmpl = EmailVerificationTemplate {
            user_name,
            verification_url,
            expiry_minutes,
        };
        let html = tmpl
            .render()
            .map_err(|e| AppError::Internal(anyhow::anyhow!("render verification email: {e}")))?;
        self.send_email(to, "Verify your email address", html).await
    }
}

// ----- askama templates (inline source) -----

/// Verification email template. Inline source keeps Rust's template
/// store separate from the Go `templates/emails/` dir that will be
/// deleted in D-REL-1.
#[derive(Template)]
#[template(
    ext = "html",
    source = r#"<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Verify your email</title>
  </head>
  <body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 560px; margin: 0 auto; padding: 24px;">
    <h2>Verify your email address</h2>
    <p>Hi {{ user_name }},</p>
    <p>Click the button below to confirm this is your email address and activate your account.</p>
    <p style="margin: 24px 0;">
      <a
        href="{{ verification_url }}"
        style="display: inline-block; padding: 12px 24px; background: #4f46e5; color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 600;"
      >Verify Email</a>
    </p>
    <p style="font-size: 12px; color: #666;">
      Or paste this URL into your browser:<br />
      <a href="{{ verification_url }}">{{ verification_url }}</a>
    </p>
    <p style="font-size: 12px; color: #666;">
      This link expires in {{ expiry_minutes }} minutes. If you didn't
      request this, you can safely ignore this email.
    </p>
  </body>
</html>"#
)]
struct EmailVerificationTemplate<'a> {
    user_name: &'a str,
    verification_url: &'a str,
    expiry_minutes: u64,
}

#[cfg(test)]
mod tests {
    use super::{EmailVerificationTemplate, Mailer};
    use askama::Template;

    #[test]
    fn verification_template_renders_all_fields() {
        let tmpl = EmailVerificationTemplate {
            user_name: "Alice",
            verification_url: "https://example.com/verify?token=abc",
            expiry_minutes: 15,
        };
        let rendered = tmpl.render().unwrap();
        assert!(rendered.contains("Hi Alice"));
        assert!(rendered.contains("https://example.com/verify?token=abc"));
        assert!(rendered.contains("15 minutes"));
    }

    #[test]
    fn noop_mailer_debug_fmt() {
        let m = Mailer::noop();
        let s = format!("{m:?}");
        assert!(s.contains("noop"));
    }
}
