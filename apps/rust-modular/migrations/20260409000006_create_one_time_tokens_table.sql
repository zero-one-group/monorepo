-- Phase D migration 6/6: one_time_tokens table + indexes.
--
-- Ported from apps/go-modular/database/migrations/00006_create_one_time_tokens_table.sql.
--
-- token_hash stays TEXT (SHA-256 hex of URL-safe token) — D-OPEN-1's
-- BYTEA harmonization only covered sessions/refresh_tokens. One-time
-- tokens hash URL-safe random strings, and TEXT storage is idiomatic
-- for human-inspectable hex.
--
-- Phase D fixes audit §9.8 (TOCTOU race on initiate/resend) via an
-- atomic `DELETE + INSERT` inside a single transaction at the service
-- layer (D-AUTH-8). The schema is unchanged — the fix is pure service
-- logic. The `UNIQUE (user_id, subject)` constraint at the bottom of
-- this migration is what makes the race visible without atomicity.

CREATE TABLE IF NOT EXISTS public.one_time_tokens (
    id           UUID NOT NULL PRIMARY KEY,
    user_id      UUID REFERENCES public.users(id) ON DELETE CASCADE,
    subject      TEXT NOT NULL,
    token_hash   TEXT NOT NULL,
    relates_to   TEXT NOT NULL,
    metadata     JSONB DEFAULT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at   TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    last_sent_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_one_time_tokens_user_id
    ON public.one_time_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_subject
    ON public.one_time_tokens (subject);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_expires_at
    ON public.one_time_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_created_at
    ON public.one_time_tokens (created_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_time_tokens_token_hash_unique
    ON public.one_time_tokens (token_hash);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_relates_to_lower
    ON public.one_time_tokens (lower(relates_to));
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_metadata_gin
    ON public.one_time_tokens USING gin (metadata);
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_time_tokens_user_id_subject
    ON public.one_time_tokens (user_id, subject);
