//! clap CLI subcommands (D-CLI-1..5).
//!
//! ```text
//! go-modular [serve]                      # default subcommand
//! go-modular migrate run                  # apply pending migrations
//! go-modular migrate create <name>        # scaffold a new migration file
//! go-modular migrate down                 # NOT SUPPORTED (no revert scripts)
//! go-modular migrate reset                # drop public schema + re-apply
//! go-modular seed                         # load dev seed data
//! go-modular generate-config [--output]   # emit .env.example
//! ```
//!
//! `migrate run` uses `sqlx::migrate!()` in-process. `migrate down`
//! errors because Phase D migrations are single-file (no separate
//! revert scripts); use `migrate reset` or restore from backup.

use std::fs;
use std::path::PathBuf;

use anyhow::{Context, Result, bail};
use clap::{Parser, Subcommand};

use crate::config::Config;
use crate::config::defaults::default_config;

/// go-modular CLI entrypoint.
#[derive(Parser, Debug)]
#[command(
    name = "go-modular",
    version,
    about = "Rust port of the Zero One Group go-modular service (Phase D)",
    long_about = None,
)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Option<Command>,
}

#[derive(Subcommand, Debug)]
pub enum Command {
    /// Start the HTTP server (default when no subcommand is given).
    Serve,

    /// Database migration management.
    Migrate {
        #[command(subcommand)]
        action: MigrateAction,
    },

    /// Load development seed data (dummy users + verified emails).
    Seed,

    /// Emit a `.env.example` file to stdout or a target path.
    GenerateConfig {
        /// Write to this file instead of stdout.
        #[arg(long, short = 'o')]
        output: Option<PathBuf>,
    },
}

#[derive(Subcommand, Debug)]
pub enum MigrateAction {
    /// Apply all pending migrations.
    Run,

    /// Create a new migration file with a timestamp prefix.
    Create {
        /// Human-readable name (`snake_case`). Example: `add_posts_table`.
        name: String,
    },

    /// Revert the latest migration.
    ///
    /// **NOT SUPPORTED in Phase D** — the ported migrations are
    /// single-file with no separate revert scripts. Use `migrate reset`
    /// (destructive) or restore the database from backup.
    Down,

    /// Drop the public schema and re-run all migrations.
    ///
    /// **Destructive.** Intended for local development only —
    /// refuses to run if `APP_MODE=production`.
    Reset,
}

/// Run the CLI, dispatching to the chosen subcommand.
pub async fn run_cli(cli: Cli) -> Result<()> {
    match cli.command.unwrap_or(Command::Serve) {
        Command::Serve => crate::server::serve().await,
        Command::Migrate { action } => match action {
            MigrateAction::Run => migrate_run().await,
            MigrateAction::Create { name } => migrate_create(&name),
            MigrateAction::Down => bail!(
                "migrate down is not supported — Phase D migrations are single-file. \
                 Use `migrate reset` (destructive) or restore from backup."
            ),
            MigrateAction::Reset => migrate_reset().await,
        },
        Command::Seed => seed_run().await,
        Command::GenerateConfig { output } => generate_config(output),
    }
}

// ----- migrate run -----

async fn migrate_run() -> Result<()> {
    let _ = dotenvy::dotenv();
    let config = Config::from_environment().context("load config")?;
    let pool = crate::database::connect_pool(&config.database).await?;
    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .context("sqlx migrate run")?;
    println!("migrate run: all pending migrations applied");
    Ok(())
}

// ----- migrate create -----

fn migrate_create(name: &str) -> Result<()> {
    let trimmed = name.trim();
    if trimmed.is_empty() {
        bail!("migration name is required");
    }
    if !trimmed
        .chars()
        .all(|c| c.is_ascii_alphanumeric() || c == '_')
    {
        bail!("migration name must be ASCII alphanumeric + underscore");
    }

    let ts = chrono::Utc::now().format("%Y%m%d%H%M%S");
    let filename = format!("{ts}_{trimmed}.sql");
    let path = PathBuf::from("apps/rust-modular/migrations").join(&filename);

    let header = format!(
        "-- Migration: {trimmed}\n\
         -- Created: {ts}\n\
         --\n\
         -- Add SQL statements below. This file is a single 'up'\n\
         -- migration — there is no separate down script. Use\n\
         -- `go-modular migrate reset` to roll back in dev.\n\n"
    );
    fs::write(&path, header).with_context(|| format!("write {}", path.display()))?;
    println!("migrate create: wrote {}", path.display());
    Ok(())
}

// ----- migrate reset -----

async fn migrate_reset() -> Result<()> {
    let _ = dotenvy::dotenv();
    let config = Config::from_environment().context("load config")?;
    if config.is_production() {
        bail!("migrate reset refuses to run in production mode");
    }
    let pool = crate::database::connect_pool(&config.database).await?;

    // Drop the public schema + sqlx's tracking table, then re-run.
    // CASCADE drops all tables, sequences, functions, triggers.
    sqlx::query("DROP SCHEMA public CASCADE")
        .execute(&pool)
        .await
        .context("drop public schema")?;
    sqlx::query("CREATE SCHEMA public")
        .execute(&pool)
        .await
        .context("recreate public schema")?;
    // Permission grants — match Postgres defaults.
    sqlx::query("GRANT ALL ON SCHEMA public TO postgres")
        .execute(&pool)
        .await
        .context("grant schema to postgres")?;
    sqlx::query("GRANT ALL ON SCHEMA public TO public")
        .execute(&pool)
        .await
        .context("grant schema to public")?;

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .context("re-run migrations")?;
    println!("migrate reset: dropped + re-applied all migrations");
    Ok(())
}

// ----- seed -----

async fn seed_run() -> Result<()> {
    let _ = dotenvy::dotenv();
    let config = Config::from_environment().context("load config")?;
    if config.is_production() {
        bail!("seed refuses to run in production mode");
    }
    let pool = crate::database::connect_pool(&config.database).await?;

    // Matches `database/seeders/user_factory.go` — create a handful
    // of verified dev users. Keep minimal: 3 users, email-verified,
    // no password (signin will fail; use `set_password` to add).
    let users = [
        ("Alice Dev", "alice@example.com", "alice"),
        ("Bob Dev", "bob@example.com", "bob"),
        ("Carol Dev", "carol@example.com", "carol"),
    ];

    for (name, email, username) in users {
        let id = uuid::Uuid::now_v7();
        sqlx::query(
            "INSERT INTO public.users \
             (id, display_name, email, username, email_verified_at) \
             VALUES ($1, $2, $3, $4, NOW()) \
             ON CONFLICT DO NOTHING",
        )
        .bind(id)
        .bind(name)
        .bind(email)
        .bind(username)
        .execute(&pool)
        .await
        .with_context(|| format!("seed user {email}"))?;
        println!("seed: {email}");
    }
    println!("seed: done");
    Ok(())
}

// ----- generate-config -----

fn generate_config(output: Option<PathBuf>) -> Result<()> {
    let cfg = default_config();
    let rendered = render_env_example(&cfg);
    match output {
        Some(path) => {
            fs::write(&path, rendered).with_context(|| format!("write {}", path.display()))?;
            println!("generate-config: wrote {}", path.display());
        }
        None => print!("{rendered}"),
    }
    Ok(())
}

fn render_env_example(cfg: &Config) -> String {
    use std::fmt::Write as _;

    let mut out = String::new();
    let _ = writeln!(out, "# Generated by `go-modular generate-config`");
    let _ = writeln!(out, "# Edit to taste; the defaults are safe for local dev.");
    let _ = writeln!(out);

    let _ = writeln!(out, "# ----- App -----");
    let _ = writeln!(out, "APP_MODE={}", cfg.app.app_mode);
    let _ = writeln!(out, "APP_BASE_URL={}", cfg.app.app_base_url);
    let _ = writeln!(out, "JWT_SECRET_KEY={}", cfg.app.jwt_secret_key);
    let _ = writeln!(out, "SERVER_HOST={}", cfg.app.server_host);
    let _ = writeln!(out, "SERVER_PORT={}", cfg.app.server_port);
    let _ = writeln!(out, "CORS_ORIGINS={}", cfg.app.cors_origins.join(","));
    let _ = writeln!(out, "CORS_MAX_AGE={}", cfg.app.cors_max_age);
    let _ = writeln!(out, "CORS_CREDENTIALS={}", cfg.app.cors_credentials);
    let _ = writeln!(out, "RATE_LIMIT_ENABLED={}", cfg.app.rate_limit_enabled);
    let _ = writeln!(out, "RATE_LIMIT_REQUESTS={}", cfg.app.rate_limit_requests);
    let _ = writeln!(
        out,
        "RATE_LIMIT_BURST_SIZE={}",
        cfg.app.rate_limit_burst_size
    );
    let _ = writeln!(out, "ENABLE_API_DOCS={}", cfg.app.enable_api_docs);

    let _ = writeln!(out, "\n# ----- Database -----");
    let _ = writeln!(out, "DATABASE_URL={}", cfg.database.database_url);
    let _ = writeln!(out, "PG_MAX_POOL_SIZE={}", cfg.database.pg_max_pool_size);
    let _ = writeln!(out, "PG_MAX_RETRIES={}", cfg.database.pg_max_retries);

    let _ = writeln!(
        out,
        "\n# ----- Mailer (lettre STARTTLS on 587, implicit TLS on 465) -----"
    );
    let _ = writeln!(out, "SMTP_HOST={}", cfg.mailer.smtp_host);
    let _ = writeln!(out, "SMTP_PORT={}", cfg.mailer.smtp_port);
    let _ = writeln!(out, "SMTP_USERNAME={}", cfg.mailer.smtp_username);
    let _ = writeln!(out, "SMTP_PASSWORD={}", cfg.mailer.smtp_password);
    let _ = writeln!(out, "SMTP_SENDER_NAME={}", cfg.mailer.smtp_sender_name);
    let _ = writeln!(out, "SMTP_SENDER_EMAIL={}", cfg.mailer.smtp_sender_email);

    let _ = writeln!(
        out,
        "\n# ----- FileStore (S3-compatible; unused by auth/user) -----"
    );
    let _ = writeln!(
        out,
        "PUBLIC_ASSETS_URL={}",
        cfg.file_store.public_assets_url
    );
    let _ = writeln!(out, "S3_ENDPOINT={}", cfg.file_store.s3_endpoint);
    let _ = writeln!(out, "S3_ACCESS_KEY={}", cfg.file_store.s3_access_key);
    let _ = writeln!(out, "S3_SECRET_KEY={}", cfg.file_store.s3_secret_key);
    let _ = writeln!(out, "S3_BUCKET_NAME={}", cfg.file_store.s3_bucket_name);
    let _ = writeln!(out, "S3_REGION={}", cfg.file_store.s3_region);
    let _ = writeln!(
        out,
        "S3_FORCE_PATH_STYLE={}",
        cfg.file_store.s3_force_path_style
    );
    let _ = writeln!(out, "S3_USE_SSL={}", cfg.file_store.s3_use_ssl);

    let _ = writeln!(out, "\n# ----- Logging -----");
    let _ = writeln!(out, "LOG_LEVEL={}", cfg.logging.log_level);
    let _ = writeln!(out, "LOG_FORMAT={}", cfg.logging.log_format);
    let _ = writeln!(out, "LOG_NO_COLOR={}", cfg.logging.log_no_color);

    let _ = writeln!(out, "\n# ----- OpenTelemetry -----");
    let _ = writeln!(out, "OTEL_SERVICE_NAME={}", cfg.otel.otel_service_name);
    let _ = writeln!(
        out,
        "OTEL_EXPORTER_OTLP_PROTOCOL={}",
        cfg.otel.otel_exporter_otlp_protocol
    );
    let _ = writeln!(
        out,
        "OTEL_EXPORTER_OTLP_ENDPOINT={}",
        cfg.otel.otel_exporter_otlp_endpoint
    );
    let _ = writeln!(
        out,
        "OTEL_EXPORTER_OTLP_HEADERS={}",
        cfg.otel.otel_exporter_otlp_headers
    );
    let _ = writeln!(
        out,
        "OTEL_ENABLE_TELEMETRY={}",
        cfg.otel.otel_enable_telemetry
    );
    let _ = writeln!(out, "OTEL_INSECURE_MODE={}", cfg.otel.otel_insecure_mode);
    let _ = writeln!(
        out,
        "OTEL_TRACING_SAMPLE_RATE={}",
        cfg.otel.otel_tracing_sample_rate
    );

    out
}
