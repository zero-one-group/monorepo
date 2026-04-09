//! Greeting service — orchestration layer.
//!
//! Mirrors `app/services/greeting.py::GreetingService`. Holds an
//! `Arc<GreetingRepoOpenAI>` and exposes a single async method
//! `greetings()` that wraps the repository response in the standard
//! `SuccessResponse` envelope.

use std::sync::Arc;

use axum::response::Response;
use tracing::Instrument;

use crate::core::exception::AppError;
use crate::core::response::success_response;
use crate::repository::openai::greeting::GreetingRepoOpenAI;

pub struct GreetingService {
    repo: Arc<GreetingRepoOpenAI>,
}

impl GreetingService {
    pub fn new(repo: Arc<GreetingRepoOpenAI>) -> Self {
        Self { repo }
    }

    /// REPLICABLE span: `service.greetings`. Must exact-match the Python
    /// `with tracer.start_as_current_span("service.greetings")` call.
    pub async fn greetings(&self) -> Result<Response, AppError> {
        let span = tracing::info_span!("service.greetings");
        async {
            tracing::info!(layer = "service", "Service layer log");
            let data = self.repo.greetings().await?;
            Ok(success_response(data))
        }
        .instrument(span)
        .await
    }
}
