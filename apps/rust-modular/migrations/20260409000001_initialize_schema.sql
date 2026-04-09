-- Phase D migration 1/6: extensions + updated_at trigger.
--
-- Ported from apps/go-modular/database/migrations/00001_initialize_schema.sql.
-- Goose directives stripped (sqlx silently no-ops on `-- +goose Up/Down`, same
-- gotcha learned in Phase C go-clean).
--
-- Delta vs Go source: the Go migrations use `DEFAULT uuidv7()` on 4 tables
-- (users, sessions, refresh_tokens, one_time_tokens). `uuidv7()` is a
-- Postgres 18+ builtin. Since our testcontainer is Postgres 16-alpine
-- (pinned in workspace deps for consistency with Phase C), we strip the
-- DB defaults from all 4 tables and generate UUIDs app-side via
-- `uuid::Uuid::now_v7()`. This is the Phase C precedent (same decision
-- applied in go-clean) and keeps compatibility with Postgres 16-alpine.

SET timezone = 'UTC';

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS plpgsql;

CREATE OR REPLACE FUNCTION fn_updated_at_value()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
