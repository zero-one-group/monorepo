-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- Create Users table -- Metadata can containg user-specific information
-- Example: SELECT * FROM users WHERE metadata ->> 'first_name' = 'Admin';
-- ============================================================================
CREATE TABLE IF NOT EXISTS public.users (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    display_name TEXT NOT NULL CHECK (char_length(display_name) > 0),
    email TEXT NOT NULL CHECK (
        char_length(email) > 3 AND
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    username TEXT UNIQUE,
    avatar_url TEXT,
    metadata JSONB DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT NULL,
    last_login_at TIMESTAMPTZ DEFAULT NULL,
    banned_at TIMESTAMPTZ DEFAULT NULL,
    ban_expires TIMESTAMPTZ DEFAULT NULL,
    ban_reason TEXT DEFAULT NULL,
    -- Username only allows alphanumeric characters and underscores, must be between 3 and 32 characters long
    CONSTRAINT chk_username_format CHECK (username IS NULL OR username ~ '^[a-zA-Z0-9_]{3,32}$')
);

-- Users table indexes and updated_at trigger
CREATE INDEX IF NOT EXISTS idx_users_display_name ON users USING gin (display_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_username ON users USING GIN (username gin_trgm_ops) WHERE username IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_email ON public.users (LOWER(email));
CREATE UNIQUE INDEX idx_users_normalized_username ON public.users (LOWER(username));
CREATE UNIQUE INDEX idx_users_normalized_email ON public.users (LOWER(email));
CREATE INDEX IF NOT EXISTS idx_users_banned_at ON public.users (banned_at) WHERE banned_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_ban_expires ON public.users (ban_expires) WHERE ban_expires IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_banned_expires ON public.users (banned_at, ban_expires) WHERE banned_at IS NOT NULL AND ban_expires IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON public.users (created_at);
CREATE INDEX IF NOT EXISTS idx_users_updated_at ON public.users (updated_at) WHERE updated_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_metadata_gin ON public.users USING GIN (metadata);
CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON public.users (last_login_at) WHERE last_login_at IS NOT NULL;
CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION fn_updated_at_value();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers, indexes, and table(s) (reverse order of creation)
DROP TRIGGER IF EXISTS trg_users_updated_at ON public.users;
DROP INDEX IF EXISTS idx_users_last_login_at;
DROP INDEX IF EXISTS idx_users_normalized_email;
DROP INDEX IF EXISTS idx_users_normalized_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_metadata_gin;
DROP INDEX IF EXISTS idx_users_updated_at;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_banned_at;
DROP INDEX IF EXISTS idx_users_ban_expires;
DROP INDEX IF EXISTS idx_users_banned_expires;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_display_name;
DROP TABLE IF EXISTS public.users;

-- +goose StatementEnd
