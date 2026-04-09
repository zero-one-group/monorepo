//! `GreetingRepoOpenAI` — the LLM call that produces 5 multilingual
//! greetings via the `OpenAI` Chat Completions API.
//!
//! Mirrors `app/repository/openai/greeting.py` exactly:
//! - Model: `gpt-4.1-mini`
//! - Temperature: `0.7`
//! - Max tokens: `150`
//! - `response_format`: `json_object` (so the model returns valid JSON)
//! - System prompt: same wording as Python
//! - User prompt: same wording as Python
//!
//! REPLICABLE span: `repository.fetch_greetings` with attributes
//! `llm.req.model`, `llm.req.temperature`, `llm.req.messages`,
//! `llm.resp.model`, `llm.usage.prompt_tokens`,
//! `llm.usage.completion_tokens`, `llm.usage.total_tokens`.
//!
//! REPLICABLE metrics:
//! - `repository_greetings_requests_total{model}` — incremented per call
//! - `repository_greetings_requests_failures_total{model, phase}` —
//!   incremented on `openai_call` or `json_parse` failures
//! - `repository_greetings_request_duration_seconds{model}` —
//!   histogram observed for the duration of the `OpenAI` call

use async_openai::Client;
use async_openai::config::OpenAIConfig;
use async_openai::types::chat::{
    ChatCompletionRequestMessage, ChatCompletionRequestSystemMessageArgs,
    ChatCompletionRequestUserMessageArgs, CreateChatCompletionRequestArgs, ResponseFormat,
};
use axum::http::StatusCode;
use serde_json::{Value, json};
use tracing::Instrument;

use crate::core::exception::AppError;
use crate::core::instrumentation::{
    greetings_request_duration, greetings_requests_failures, greetings_requests_total,
};

const MODEL: &str = "gpt-4.1-mini";
const TEMPERATURE: f32 = 0.7;
const MAX_TOKENS: u32 = 150;

const SYSTEM_PROMPT: &str = "You are a helpful assistant that provides greetings in different languages. Always respond in valid JSON format.";
const USER_PROMPT: &str = r#"Give me greetings in exactly 5 different languages. Return only a JSON object with this exact structure: {"greetings": [{"language": "name of language", "greeting": "greeting in that language"}]}. Include no explanations or additional text."#;

pub struct GreetingRepoOpenAI {
    client: Client<OpenAIConfig>,
}

impl GreetingRepoOpenAI {
    pub fn new(client: Client<OpenAIConfig>) -> Self {
        Self { client }
    }

    /// Fetch greetings. Returns the parsed JSON wrapped in a
    /// `{"response": {...}}` object, matching the Python repo's
    /// `return {"response": greetings_data}`.
    pub async fn greetings(&self) -> Result<Value, AppError> {
        greetings_requests_total().with_label_values(&[MODEL]).inc();

        tracing::info!(
            layer = "repository",
            "Fetching greetings in 5 languages from OpenAI"
        );

        let timer = greetings_request_duration()
            .with_label_values(&[MODEL])
            .start_timer();

        let span = tracing::info_span!(
            "repository.fetch_greetings",
            llm.req.model = MODEL,
            llm.req.temperature = TEMPERATURE,
            // The Python code attaches the messages array as an attribute.
            // We attach the joined content for parity (axum/tracing don't
            // support array-valued attributes natively).
            llm.req.messages = format!("{SYSTEM_PROMPT} || {USER_PROMPT}").as_str()
        );

        let result = self.call_openai_and_parse().instrument(span).await;

        timer.observe_duration();

        result
    }

    async fn call_openai_and_parse(&self) -> Result<Value, AppError> {
        let messages: Vec<ChatCompletionRequestMessage> = vec![
            ChatCompletionRequestSystemMessageArgs::default()
                .content(SYSTEM_PROMPT)
                .build()
                .map_err(|e| {
                    AppError::new(
                        format!("OpenAI request build failed: {e}"),
                        StatusCode::INTERNAL_SERVER_ERROR,
                        "OPENAI_BUILD_ERROR",
                    )
                })?
                .into(),
            ChatCompletionRequestUserMessageArgs::default()
                .content(USER_PROMPT)
                .build()
                .map_err(|e| {
                    AppError::new(
                        format!("OpenAI request build failed: {e}"),
                        StatusCode::INTERNAL_SERVER_ERROR,
                        "OPENAI_BUILD_ERROR",
                    )
                })?
                .into(),
        ];

        let request = CreateChatCompletionRequestArgs::default()
            .model(MODEL)
            .messages(messages)
            .temperature(TEMPERATURE)
            .max_tokens(MAX_TOKENS)
            .response_format(ResponseFormat::JsonObject)
            .build()
            .map_err(|e| {
                AppError::new(
                    format!("OpenAI request build failed: {e}"),
                    StatusCode::INTERNAL_SERVER_ERROR,
                    "OPENAI_BUILD_ERROR",
                )
            })?;

        let response = match self.client.chat().create(request).await {
            Ok(r) => r,
            Err(e) => {
                greetings_requests_failures()
                    .with_label_values(&[MODEL, "openai_call"])
                    .inc();
                return Err(AppError::new(
                    format!("OpenAI unavailable: {e}"),
                    StatusCode::SERVICE_UNAVAILABLE,
                    "OPENAI_ERROR",
                ));
            }
        };

        let resp_model = response.model.clone();
        if let Some(usage) = response.usage.as_ref() {
            tracing::Span::current().record("llm.resp.model", resp_model.as_str());
            tracing::Span::current().record("llm.usage.prompt_tokens", usage.prompt_tokens);
            tracing::Span::current().record("llm.usage.completion_tokens", usage.completion_tokens);
            tracing::Span::current().record("llm.usage.total_tokens", usage.total_tokens);
        }

        let content = response
            .choices
            .first()
            .and_then(|c| c.message.content.clone())
            .unwrap_or_default();

        tracing::info!(
            layer = "repository",
            "Successfully received greetings from OpenAI"
        );

        let parsed: Value = match serde_json::from_str(&content) {
            Ok(v) => v,
            Err(e) => {
                greetings_requests_failures()
                    .with_label_values(&[MODEL, "json_parse"])
                    .inc();
                return Err(AppError::new(
                    format!("Invalid JSON from OpenAI: {e}"),
                    StatusCode::BAD_GATEWAY,
                    "PARSE_ERROR",
                ));
            }
        };

        let language_count = parsed
            .get("greetings")
            .and_then(|g| g.as_array())
            .map_or(0, std::vec::Vec::len);
        tracing::info!(
            layer = "repository",
            "Successfully parsed greetings in {language_count} languages"
        );

        Ok(json!({ "response": parsed }))
    }
}
