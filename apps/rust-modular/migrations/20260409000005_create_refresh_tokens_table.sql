-- Phase D migration 5/6: refresh_tokens table + indexes.
--
-- Ported from apps/go-modular/database/migrations/00005_create_refresh_tokens_table.sql.
--
-- SEMANTIC DRIFT vs Go source (D-OPEN-1 decision, 2026-04-09):
--   token_hash BYTEA is unchanged in type but NOW holds the raw 32-byte
--   SHA-256 digest, not the ASCII bytes of a 64-char hex string (which
--   is what the Go code stored). The Rust signin service computes
--   `Sha256::digest(jwt.as_bytes())` directly and passes the `[u8; 32]`
--   to sqlx — no `hex::encode` intermediate.
--
-- The `update` method is DELETED in the Rust repository per D-OPEN-5:
-- the 4 naive CRUD endpoints are gone, replaced with ONE rotation
-- endpoint (POST /api/v1/auth/token/refresh). This migration's table
-- is unchanged; only the access pattern differs.

CREATE TABLE IF NOT EXISTS public.refresh_tokens (
    id         UUID NOT NULL PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES public.sessions(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ DEFAULT NULL,
    revoked_by UUID REFERENCES public.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id
    ON public.refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_session_id
    ON public.refresh_tokens (session_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at
    ON public.refresh_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at
    ON public.refresh_tokens (revoked_at) WHERE revoked_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_created_at
    ON public.refresh_tokens (created_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_ip_address
    ON public.refresh_tokens (ip_address) WHERE ip_address IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_agent
    ON public.refresh_tokens (user_agent) WHERE user_agent IS NOT NULL;
