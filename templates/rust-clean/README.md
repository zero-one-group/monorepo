# {{ package_name | kebab_case }}

{{ package_description }}

Rust/axum clean architecture REST API with JWT authentication, bcrypt password hashing,
rate limiting, and PostgreSQL persistence via sqlx.

## Prerequisites

- Rust 1.90+
- [Moonrepo](https://moonrepo.dev/docs/getting-started/installation)
- PostgreSQL 16+
- Docker and Docker Compose (optional)

## Quick Start

1. Generate a new project from this template:

    ```bash
    moon generate rust-clean
    ```

2. Set up environment variables:

    ```bash
    cp .env.example .env
    # Edit .env — set DATABASE_URL, JWT_SECRET_KEY, etc.
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
| `POST` | `/auth/register` | Register a new user |
| `POST` | `/auth/login` | Login and receive JWT |
| `GET` | `/users` | List users (authenticated) |
| `GET` | `/users/:id` | Get user by ID (authenticated) |
| `PUT` | `/users/:id` | Update user (authenticated) |

## Architecture

```
src/
├── config/       # Env config loader (figment)
├── domain/       # Domain types: User, Auth, Error, Response
├── repository/   # Postgres data access layer
├── rest/         # Axum handlers + JWT middleware
├── service/      # Business logic
└── utils/        # JWT helpers, password hashing
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

Raise `AppError` from `domain/error.rs` — the global handler converts it automatically.
