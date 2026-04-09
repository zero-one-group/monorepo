//! `templates-cli` library crate.
//!
//! Contains the pure-logic pieces of the template-build pipeline so they
//! can be unit-tested without invoking the binary.
//!
//! The public surface:
//! - [`Cli`] — top-level clap command
//! - [`run`] — entrypoint the binary calls with a parsed `Cli`
//! - [`common`] — shared helpers (placeholder replacement, file renaming,
//!   JSON surgery, SHA-256 manifesting)
//! - [`commands::build`] — `templates-cli build` (= build-templates.sh)
//! - [`commands::makezip`] — `templates-cli makezip` (= makezip.sh)
//! - [`builders`] — one module per `builder/*.sh` script

pub mod builders;
pub mod commands;
pub mod common;

use std::path::PathBuf;

use anyhow::Result;
use clap::{Parser, Subcommand};

/// Top-level CLI.
#[derive(Debug, Parser)]
#[command(name = "templates-cli", version, about)]
pub struct Cli {
    /// Path to the monorepo root. Defaults to the current working directory.
    #[arg(long, global = true)]
    pub root: Option<PathBuf>,

    #[command(subcommand)]
    pub command: Command,
}

#[derive(Debug, Subcommand)]
pub enum Command {
    /// Run the full build + zip pipeline (build-templates.sh + makezip.sh).
    All,

    /// Run only the build phase (build-templates.sh).
    ///
    /// Copies every app in the source→target map into `templates/`, strips
    /// runtime artifacts, then invokes every registered builder module.
    Build,

    /// Run only the zip phase (makezip.sh).
    ///
    /// Zips every `templates/<name>/` subdirectory into
    /// `docsite/static/templates/<name>.zip` and regenerates
    /// `docsite/static/templates.json`.
    Makezip,

    /// Run one builder module by name.
    ///
    /// Names match the original `builder/*.sh` scripts:
    /// `astro`, `expo-app`, `fastapi-ai`, `go-clean`, `go-modular`,
    /// `nextjs-app`, `react-app`, `react-ssr`, `shared-ui`, `strapi-cms`,
    /// `tanstack-start`.
    Builder {
        /// Builder to run.
        name: String,
    },
}

/// Entrypoint invoked from `main.rs`.
pub fn run(cli: Cli) -> Result<()> {
    let root = match cli.root {
        Some(p) => p,
        None => std::env::current_dir()?,
    };
    tracing::debug!(?root, "resolved monorepo root");

    match cli.command {
        Command::All => {
            commands::build::run(&root)?;
            commands::makezip::run(&root)?;
        }
        Command::Build => {
            commands::build::run(&root)?;
        }
        Command::Makezip => {
            commands::makezip::run(&root)?;
        }
        Command::Builder { name } => {
            builders::run_builder(&name, &root)?;
        }
    }
    Ok(())
}
