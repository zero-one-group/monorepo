# {{ package_name | kebab_case }}

{{ package_description }}

Rust/axum AI service with OpenAI-backed greeting endpoints, OpenTelemetry tracing,
Prometheus metrics, and PostgreSQL persistence via sqlx.

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
        "apps/{{ package_name | kebab_case }}",
    ]
    ```

4. Run migrations and start:

    ```bash
    moon run {{ package_name | kebab_case }}:migrate
    moon run {{ package_name | kebab_case }}:dev
    ```

## Available Commands

| Command | Description |
|---|---|
| `moon run {{ package_name | kebab_case }}:dev` | Run in development mode |
| `moon run {{ package_name | kebab_case }}:start` | Run release build |
| `moon run {{ package_name | kebab_case }}:build` | Build debug binary |
| `moon run {{ package_name | kebab_case }}:build-release` | Build release binary |
| `moon run {{ package_name | kebab_case }}:test` | Run all tests |
| `moon run {{ package_name | kebab_case }}:lint` | Run clippy lints |
| `moon run {{ package_name | kebab_case }}:migrate` | Apply pending migrations |
| `moon run {{ package_name | kebab_case }}:migrate-create -- "name"` | Create a new migration file |
| `moon run {{ package_name | kebab_case }}:migrate-down` | Revert last migration |
| `moon run {{ package_name | kebab_case }}:migrate-reset` | Reset database (destroys all data) |
| `moon run {{ package_name | kebab_case }}:check-in-dance` | Build + migrate (first-time setup) |

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

## Development

Walk through the "Greeting API" example:

1. Start at the router: `src/router/openai.rs`
2. Inspect the service layer: `src/services/greeting.rs`
3. Examine the repository layer: `src/repository/openai/greeting.rs`

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
