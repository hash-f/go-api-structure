CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Optional: Add an extension for gen_random_uuid() if not already enabled
-- CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- Note: gen_random_uuid() is available in PostgreSQL 13+.
-- For older versions, you might need pgcrypto's uuid_generate_v4().
