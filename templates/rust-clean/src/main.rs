use std::process::ExitCode;

use {{ package_name | snake_case }}::serve;

#[tokio::main]
async fn main() -> ExitCode {
    if let Err(err) = serve().await {
        eprintln!("{{ package_name | kebab_case }}: fatal error: {err:#}");
        return ExitCode::FAILURE;
    }
    ExitCode::SUCCESS
}
