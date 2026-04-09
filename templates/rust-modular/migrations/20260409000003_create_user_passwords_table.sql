-- Phase D migration 3/6: user_passwords table + updated_at trigger.
--
-- Ported from apps/go-modular/database/migrations/00003_create_user_passwords_table.sql.
-- BYTEA password_hash stores a PHC-encoded argon2id string. Phase D bumps
-- the argon2id salt length from Go's 8 bytes to OWASP's 16 bytes per
-- D-OPEN-2; the PHC format encodes salt length per-hash so any Go-era
-- hashes remain verifiable via `PasswordHash::new(&phc_string)`.

CREATE TABLE IF NOT EXISTS public.user_passwords (
    user_id       UUID NOT NULL PRIMARY KEY REFERENCES public.users(id) ON DELETE CASCADE,
    password_hash BYTEA NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT NULL,
    CONSTRAINT user_passwords_one_per_user UNIQUE (user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_passwords_created_at
    ON public.user_passwords (created_at);
CREATE INDEX IF NOT EXISTS idx_user_passwords_updated_at
    ON public.user_passwords (updated_at) WHERE updated_at IS NOT NULL;

CREATE TRIGGER trg_user_passwords_updated_at
    BEFORE UPDATE ON public.user_passwords
    FOR EACH ROW EXECUTE FUNCTION fn_updated_at_value();
