-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- Create user_passwords table and indexes
-- SELECT user_id, convert_from(password_hash, 'UTF8') AS password_hash, created_at, updated_at FROM public.user_passwords;
-- ============================================================================
CREATE TABLE IF NOT EXISTS public.user_passwords (
    user_id UUID NOT NULL PRIMARY KEY REFERENCES public.users(id) ON DELETE CASCADE,
    password_hash BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT NULL,
    CONSTRAINT user_passwords_one_per_user UNIQUE (user_id) -- Ensure only one password per user
);

-- Users passwords table indexes and updated_at trigger
CREATE INDEX IF NOT EXISTS idx_user_passwords_created_at ON public.user_passwords (created_at);
CREATE INDEX IF NOT EXISTS idx_user_passwords_updated_at ON public.user_passwords (updated_at) WHERE updated_at IS NOT NULL;
CREATE TRIGGER trg_user_passwords_updated_at BEFORE UPDATE ON public.user_passwords FOR EACH ROW EXECUTE FUNCTION fn_updated_at_value();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers, indexes, and table(s) (reverse order of creation)
DROP TRIGGER IF EXISTS trg_user_passwords_updated_at ON public.user_passwords;
DROP INDEX IF EXISTS idx_user_passwords_updated_at;
DROP INDEX IF EXISTS idx_user_passwords_created_at;
DROP TABLE IF EXISTS public.user_passwords;

-- +goose StatementEnd
