-- Translated from the original goose migration at this same path on
-- the `legacy/go-clean-original` branch (pre-Phase-C state):
-- `-- +goose Up` directives and the `CREATE EXTENSION "uuid-ossp"`
-- preamble have been removed because sqlx's migrator does NOT parse
-- goose's custom comment syntax and silently treats the file as
-- zero-statement no-op. Discovered during Phase C golden test setup.
--
-- One deliberate schema deviation: the Go migration uses
-- `uuid_generate_v4()` from the `uuid-ossp` extension as the PK
-- default. That extension is NOT loaded by default on the postgres
-- image that `testcontainers-modules::Postgres::default()` spins up.
-- `gen_random_uuid()` is a core PostgreSQL function since version 13
-- (no extension required), functionally equivalent to
-- `uuid_generate_v4()`. We use it here so the same migration file
-- works against both the shared dev Postgres (which has uuid-ossp
-- available if needed) and fresh testcontainer Postgres images.

CREATE TABLE IF NOT EXISTS users (
    id         UUID                      PRIMARY KEY,
    name       VARCHAR                   NOT NULL UNIQUE,
    email      VARCHAR                   NOT NULL UNIQUE,
    password   TEXT                      NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE  DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE  DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- Note: no DB-side default for `id`. The Rust application layer
-- generates a fresh v4 UUID before every INSERT and passes it
-- explicitly, making this migration portable across postgres
-- versions (`gen_random_uuid()` requires postgres 13+ in core) and
-- across test images that may not ship the `uuid-ossp` extension.
