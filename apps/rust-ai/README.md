# rust-ai

Rust/axum port of the AI service. Provides OpenAI-backed greeting endpoints with
OpenTelemetry tracing, Prometheus metrics, and PostgreSQL persistence.

## Prerequisites

- Rust 1.90+
- [Moonrepo](https://moonrepo.dev/docs/getting-started/installation)
- PostgreSQL 16+
- Docker and Docker Compose (optional)

## Quick Start

1. Generate a new project from this template:

    ```bash
    moon generate rust-ai
    ```

2. Set up environment variables:

    ```bash
    cp .env.example .env
    # Edit .env — set DATABASE_URL and OPENAI_API_KEY
    ```

3. Add the new crate to the workspace (`Cargo.toml`):

    ```toml
    members = [
        # ...
        "apps/your-app-name",
    ]
    ```

4. Run migrations and start:

    ```bash
    moon run rust-ai:migrate
    moon run rust-ai:dev
    ```

## Available Commands

| Command | Description |
|---|---|
| `moon run rust-ai:dev` | Run in development mode |
| `moon run rust-ai:start` | Run release build |
| `moon run rust-ai:build` | Build debug binary |
| `moon run rust-ai:build-release` | Build release binary |
| `moon run rust-ai:test` | Run all tests |
| `moon run rust-ai:lint` | Run clippy lints |
| `moon run rust-ai:migrate` | Apply pending migrations |
| `moon run rust-ai:migrate-create -- "name"` | Create a new migration file |
| `moon run rust-ai:migrate-down` | Revert last migration |
| `moon run rust-ai:migrate-reset` | Reset database (destroys all data) |
| `moon run rust-ai:check-in-dance` | Build + migrate (first-time setup) |

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Root health check |
| `GET` | `/health-check` | Liveness probe |
| `GET` | `/openai/greetings` | AI greeting via OpenAI |

## Architecture

```
src/
├── core/         # Config, DB pool, error types, logging, OTel, metrics
├── model/        # Domain models
├── repository/   # Data access (OpenAI, Postgres)
├── router/       # Axum route handlers
└── services/     # Business logic
```

## Error Handling

All errors return a structured JSON envelope:

```json
{
  "success": false,
  "message": "description",
  "error_code": "CODE",
  "data": null
}
```

Raise `AppError` from `core/exception.rs` at any layer — the global handler converts it automatically.

## Observability

- **Tracing**: OTel spans emitted in production (`APP_ENVIRONMENT=production`)
- **Metrics**: Prometheus registry exposed at `/metrics`
- Structural gaps vs the Python original documented in `src/core/instrumentation.rs`
