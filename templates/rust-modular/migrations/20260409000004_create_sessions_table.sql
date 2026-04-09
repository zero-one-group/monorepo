-- Phase D migration 4/6: sessions table + indexes.
--
-- Ported from apps/go-modular/database/migrations/00004_create_sessions_table.sql.
--
-- DELIBERATE DRIFT vs Go source (D-OPEN-1 decision, 2026-04-09):
--   Go:  token_hash TEXT NOT NULL UNIQUE  (ASCII bytes of 64-char hex digest)
--   Rust: token_hash BYTEA NOT NULL UNIQUE (raw 32-byte SHA-256 digest)
--
-- Rationale: fresh data means zero migration cost, and harmonizing both
-- sessions.token_hash and refresh_tokens.token_hash to the same raw
-- BYTEA representation eliminates audit §9.17's schema inconsistency.

CREATE TABLE IF NOT EXISTS public.sessions (
    id                 UUID NOT NULL PRIMARY KEY,
    user_id            UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    token_hash         BYTEA NOT NULL UNIQUE,
    user_agent         TEXT,
    device_name        TEXT,
    device_fingerprint TEXT,
    ip_address         INET,
    expires_at         TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    refreshed_at       TIMESTAMPTZ DEFAULT NULL,
    revoked_at         TIMESTAMPTZ DEFAULT NULL,
    revoked_by         UUID REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id
    ON public.sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at
    ON public.sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id_expires_at
    ON public.sessions (user_id, expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_ip_address
    ON public.sessions (ip_address) WHERE ip_address IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_device_fingerprint
    ON public.sessions (device_fingerprint);
