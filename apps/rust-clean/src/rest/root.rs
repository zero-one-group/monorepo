//! Root `/` health endpoint. Mirrors the `GET /` handler in `main.go`.

use axum::Json;

use crate::domain::response::Response;

pub async fn health() -> Json<Response> {
    Json(Response::new(200, "All is well!"))
}
