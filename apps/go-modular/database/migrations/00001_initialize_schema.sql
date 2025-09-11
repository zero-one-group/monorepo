-- +goose Up
-- +goose StatementBegin
SET timezone = 'UTC';

-- ============================================================================
-- Register some PosgreSQL extensions
-- ============================================================================
-- CREATE EXTENSION IF NOT EXISTS pg_stat_statements; -- Track planning and execution statistics of all SQL statements executed
-- CREATE EXTENSION IF NOT EXISTS citext;    -- Case-insensitive text type (slower than varchar or text columns)
-- CREATE EXTENSION IF NOT EXISTS hstore;    -- Key-Value pairs for storing unstructured data
CREATE EXTENSION IF NOT EXISTS pg_trgm;   -- Text similarity measurement and index searching based on trigrams
CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- Cryptographic functions for hashing and encryption
CREATE EXTENSION IF NOT EXISTS plpgsql;   -- PL/pgSQL procedural language

-- ============================================================================
-- Create auto-update function, fill updated_at column automatically.
-- CURRENT_TIMESTAMP similar to timezone('utc'::text, now())::timestamptz
-- ============================================================================
CREATE OR REPLACE FUNCTION fn_updated_at_value()
RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $$
LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop custom functions but keep the extensions.
DROP FUNCTION IF EXISTS fn_updated_at_value();

-- +goose StatementEnd
