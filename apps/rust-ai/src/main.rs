use std::process::ExitCode;

use rust_ai::serve;

#[tokio::main]
async fn main() -> ExitCode {
    if let Err(err) = serve().await {
        eprintln!("fastapi-ai: fatal error: {err:#}");
        return ExitCode::FAILURE;
    }
    ExitCode::SUCCESS
}
