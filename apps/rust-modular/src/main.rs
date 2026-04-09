//! go-modular entrypoint — delegates to `cli::run_cli`.

use std::process::ExitCode;

use clap::Parser;
use rust_modular::cli::{Cli, run_cli};

#[tokio::main]
async fn main() -> ExitCode {
    let cli = Cli::parse();
    if let Err(err) = run_cli(cli).await {
        eprintln!("go-modular: error: {err:#}");
        return ExitCode::FAILURE;
    }
    ExitCode::SUCCESS
}
