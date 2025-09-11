-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- Create auth refresh token table and indexes
-- ============================================================================
CREATE TABLE IF NOT EXISTS public.refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES public.sessions(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ DEFAULT NULL,
    revoked_by UUID REFERENCES public.users(id) ON DELETE SET NULL
);

-- Refresh tokens table indexes
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON public.refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_session_id ON public.refresh_tokens (session_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON public.refresh_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at ON public.refresh_tokens (revoked_at) WHERE revoked_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_created_at ON public.refresh_tokens (created_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_ip_address ON public.refresh_tokens (ip_address) WHERE ip_address IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_agent ON public.refresh_tokens (user_agent) WHERE user_agent IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON public.refresh_tokens (token_hash);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes, and table(s) (reverse order of creation)
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_session_id;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_revoked_at;
DROP INDEX IF EXISTS idx_refresh_tokens_created_at;
DROP INDEX IF EXISTS idx_refresh_tokens_ip_address;
DROP INDEX IF EXISTS idx_refresh_tokens_user_agent;
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;
DROP TABLE IF EXISTS public.refresh_tokens;

-- +goose StatementEnd
