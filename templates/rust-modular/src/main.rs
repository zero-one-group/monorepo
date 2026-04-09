//! {{ package_name | kebab_case }} entrypoint — delegates to `cli::run_cli`.

use std::process::ExitCode;

use clap::Parser;
use {{ package_name | snake_case }}::cli::{Cli, run_cli};

#[tokio::main]
async fn main() -> ExitCode {
    let cli = Cli::parse();
    if let Err(err) = run_cli(cli).await {
        eprintln!("{{ package_name | kebab_case }}: error: {err:#}");
        return ExitCode::FAILURE;
    }
    ExitCode::SUCCESS
}
