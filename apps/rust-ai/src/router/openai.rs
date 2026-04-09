//! `/openai/greetings` route.
//!
//! Mirrors `app/router/openai.py`. Returns the `SuccessResponse` envelope
//! produced by `services::greeting::GreetingService::greetings()`.
//!
//! Span: `route.greetings` with attributes `http.method=GET` and
//! `http.route=/greetings`. This is REPLICABLE — must exact-match.

use axum::extract::State;
use axum::response::Response;
use axum::routing::get;
use tracing::Instrument;

use crate::AppState;
use crate::core::exception::AppError;

pub fn router() -> axum::Router<AppState> {
    axum::Router::new().route("/openai/greetings", get(greetings))
}

async fn greetings(State(state): State<AppState>) -> Result<Response, AppError> {
    let span = tracing::info_span!(
        "route.greetings",
        http.method = "GET",
        http.route = "/greetings"
    );
    state.greeting_service.greetings().instrument(span).await
}
