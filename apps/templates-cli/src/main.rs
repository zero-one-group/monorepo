use std::process::ExitCode;

use clap::Parser;
use templates_cli::{Cli, run};
use tracing_subscriber::EnvFilter;

fn main() -> ExitCode {
    // Route tracing to stderr in the same human-readable shape the bash
    // scripts produced (one line per action). `RUST_LOG=debug` unlocks
    // deeper diagnostics during implementation and tests.
    tracing_subscriber::fmt()
        .with_env_filter(
            EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info")),
        )
        .with_target(false)
        .without_time()
        .with_writer(std::io::stderr)
        .init();

    let cli = Cli::parse();
    match run(cli) {
        Ok(()) => ExitCode::SUCCESS,
        Err(err) => {
            tracing::error!("{err:#}");
            ExitCode::FAILURE
        }
    }
}
