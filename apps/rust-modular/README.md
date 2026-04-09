# rust-modular

NestJS-style modular Rust/axum service with full auth system (14 endpoints),
user management (5 endpoints), JWT, argon2 password hashing, SMTP mailer,
rate limiting, OpenAPI docs, and PostgreSQL via sqlx.

## Prerequisites

- Rust 1.90+
- [Moonrepo](https://moonrepo.dev/docs/getting-started/installation)
- PostgreSQL 16+
- Docker and Docker Compose (optional, MailHog for local SMTP)

## Quick Start

1. Generate a new project from this template:

    ```bash
    moon generate rust-modular
    ```

2. Set up environment variables:

    ```bash
    cp .env.example .env
    # Edit .env — set DATABASE_URL, JWT_SECRET_KEY, SMTP_* etc.
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
    moon run rust-modular:migrate
    moon run rust-modular:dev
    ```

## Available Commands

| Command | Description |
|---|---|
| `moon run rust-modular:dev` | Run in development mode |
| `moon run rust-modular:start` | Run release build |
| `moon run rust-modular:build` | Build debug binary |
| `moon run rust-modular:build-release` | Build release binary |
| `moon run rust-modular:test` | Run all tests |
| `moon run rust-modular:lint` | Run clippy lints |
| `moon run rust-modular:migrate` | Apply pending migrations |
| `moon run rust-modular:migrate-create -- "name"` | Create a new migration file |
| `moon run rust-modular:migrate-down` | Revert last migration |
| `moon run rust-modular:migrate-reset` | Reset database (destroys all data) |
| `moon run rust-modular:check-in-dance` | Build + migrate (first-time setup) |

## CLI

The binary also exposes a CLI:

```bash
cargo run -p rust-modular -- serve           # start HTTP server (default)
cargo run -p rust-modular -- migrate run     # apply pending migrations
cargo run -p rust-modular -- migrate create "name"  # scaffold migration
cargo run -p rust-modular -- migrate reset   # drop + re-apply all
cargo run -p rust-modular -- seed            # load dev seed data
cargo run -p rust-modular -- generate-config # emit .env.example to stdout
```

## Endpoints

### Auth (14 endpoints)

| Method | Path | Description |
|---|---|---|
| `POST` | `/auth/register` | Register new user |
| `POST` | `/auth/login` | Login, returns access + refresh tokens |
| `POST` | `/auth/refresh` | Rotate refresh token |
| `DELETE` | `/auth/logout` | Revoke session |
| `GET` | `/auth/profile` | Get own profile |
| `POST` | `/auth/forgot-password` | Send reset email |
| `POST` | `/auth/reset-password` | Apply reset token |
| `POST` | `/auth/verify-email` | Verify email address |
| `POST` | `/auth/resend-verification` | Resend verification email |
| `PUT` | `/auth/change-password` | Change password |
| `DELETE` | `/auth/sessions` | Revoke all sessions |
| `GET` | `/auth/sessions` | List active sessions |
| `DELETE` | `/auth/sessions/:id` | Revoke specific session |
| `GET` | `/auth/me` | Alias for profile |

### User (5 endpoints)

| Method | Path | Description |
|---|---|---|
| `GET` | `/users` | List users |
| `GET` | `/users/:id` | Get user by ID |
| `PUT` | `/users/:id` | Update user |
| `DELETE` | `/users/:id` | Delete user |
| `GET` | `/users/me` | Get own user record |

### Other

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthz` | Liveness probe |
| `GET` | `/api-docs` | Swagger UI (when `ENABLE_API_DOCS=true`) |

## Architecture

```
src/
├── apputils/     # JWT, password hashing, validation, token generation
├── cli/          # Clap CLI (serve / migrate / seed / generate-config)
├── config/       # Figment config loader (6 sections, 39 fields)
├── database/     # PgPool with retry
├── domain/       # Shared error types and response envelopes
├── mailer/       # Lettre SMTP mailer with askama email templates
├── middleware/   # Tower middleware stack (rate limit, request-id, etc.)
├── modules/
│   ├── auth/     # handler / service / repository / models / schema
│   └── user/     # handler / service / repository / models
├── observer/     # Tracing subscriber + OTel wiring
├── openapi/      # Utoipa OpenAPI spec + Swagger UI
└── server/       # Axum router builder + serve loop
```
