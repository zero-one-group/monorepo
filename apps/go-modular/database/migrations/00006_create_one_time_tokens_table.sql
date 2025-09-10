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
    metadata JSONB DEFAULT NULL, -- additional information if needed
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > CURRENT_TIMESTAMP),
    last_sent_at TIMESTAMPTZ DEFAULT NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_user_id ON public.one_time_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_subject ON public.one_time_tokens (subject);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_expires_at ON public.one_time_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_created_at ON public.one_time_tokens (created_at);

-- Use a btree unique index on token_hash for fast exact lookups and to prevent duplicates.
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_time_tokens_token_hash_unique ON public.one_time_tokens (token_hash);

-- Index on lower(relates_to) to support case-insensitive lookups for emails/handles.
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_relates_to_lower ON public.one_time_tokens (lower(relates_to));

-- GIN index for metadata JSONB queries
CREATE INDEX IF NOT EXISTS idx_one_time_tokens_metadata_gin ON public.one_time_tokens USING GIN (metadata);

-- Unique constraint to ensure a user has at most one active token per subject (e.g. email_verification)
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_time_tokens_user_id_subject ON public.one_time_tokens (user_id, subject);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes and table (reverse order)
DROP INDEX IF EXISTS idx_one_time_tokens_user_id_subject;
DROP INDEX IF EXISTS idx_one_time_tokens_metadata_gin;
DROP INDEX IF EXISTS idx_one_time_tokens_relates_to_lower;
DROP INDEX IF EXISTS idx_one_time_tokens_token_hash_unique;
DROP INDEX IF EXISTS idx_one_time_tokens_created_at;
DROP INDEX IF EXISTS idx_one_time_tokens_expires_at;
DROP INDEX IF EXISTS idx_one_time_tokens_subject;
DROP INDEX IF EXISTS idx_one_time_tokens_user_id;

DROP TABLE IF EXISTS public.one_time_tokens;

-- +goose StatementEnd
