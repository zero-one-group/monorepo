use std::process::ExitCode;

use rust_clean::serve;

#[tokio::main]
async fn main() -> ExitCode {
    if let Err(err) = serve().await {
        eprintln!("go-clean: fatal error: {err:#}");
        return ExitCode::FAILURE;
    }
    ExitCode::SUCCESS
}
