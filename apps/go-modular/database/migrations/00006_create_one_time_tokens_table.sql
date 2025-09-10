-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- Create One Time Token table and indexes
-- Used for email verification, password reset, signin with OTP, etc.
-- ============================================================================
CREATE TABLE IF NOT EXISTS public.one_time_tokens (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID REFERENCES public.users(id) ON DELETE CASCADE,
    subject TEXT NOT NULL, -- email_otp, mfa_otp, email_verification, reauthentication, etc.
    token_hash TEXT NOT NULL, -- hashed token (SHA256)
    relates_to TEXT NOT NULL, -- value can be email, phone, etc.
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    last_sent_at TIMESTAMPTZ DEFAULT NULL
);

-- One Time Token indexes
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_user_id ON public.one_time_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_subject ON public.one_time_tokens (subject);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_expires_at ON public.one_time_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_token_hash ON public.one_time_tokens USING HASH (token_hash);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_relates_to ON public.one_time_tokens USING HASH (relates_to);
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_time_tokens_user_id_subject ON public.one_time_tokens (user_id, subject);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes, and table(s) (reverse order of creation)
DROP INDEX IF EXISTS idx_one_time_tokens_created_at;
DROP INDEX IF EXISTS idx_one_time_tokens_expires_at;
DROP INDEX IF EXISTS idx_one_time_tokens_user_id_subject;
DROP INDEX IF EXISTS idx_one_time_tokens_relates_to;
DROP INDEX IF EXISTS idx_one_time_tokens_token_hash;
DROP INDEX IF EXISTS idx_one_time_tokens_token;
DROP INDEX IF EXISTS idx_one_time_tokens_user_id;
DROP TABLE IF EXISTS public.one_time_tokens;

-- +goose StatementEnd
