-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- Create auth sessions table and indexes
-- ============================================================================
CREATE TABLE IF NOT EXISTS public.sessions (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    device_name TEXT,
    device_fingerprint TEXT,
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    refreshed_at TIMESTAMPTZ DEFAULT NULL,
    revoked_at TIMESTAMPTZ DEFAULT NULL,
    revoked_by UUID REFERENCES public.users(id) ON DELETE SET NULL
);

-- Sessions table indexes
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON public.sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON public.sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id_expires_at ON public.sessions (user_id, expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_ip_address ON public.sessions (ip_address) WHERE ip_address IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_device_fingerprint ON public.sessions (device_fingerprint);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_token_hash ON public.sessions (token_hash);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes, and table(s) (reverse order of creation)
DROP INDEX IF EXISTS idx_sessions_ip_address;
DROP INDEX IF EXISTS idx_sessions_user_id_expires_at;
DROP INDEX IF EXISTS idx_sessions_device_fingerprint;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_token_hash;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP TABLE IF EXISTS public.sessions;

-- +goose StatementEnd
