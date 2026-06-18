-- 000001_create_users_table.up.sql
-- Creates the main users table with audit columns and trigger

CREATE TABLE IF NOT EXISTS tbl_users (
    id BIGSERIAL PRIMARY KEY,
    user_name VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100),
    role_name VARCHAR(50) NOT NULL DEFAULT 'SuperAdmin',
    role_id INTEGER NOT NULL DEFAULT 0,
    is_admin BOOLEAN DEFAULT false,
    login_session VARCHAR(255),
    last_login TIMESTAMP,
    currency_id INTEGER,
    language_id INTEGER,
    status_id INTEGER DEFAULT 1,
    "order" INTEGER,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_by BIGINT,
    updated_at TIMESTAMP,
    deleted_by BIGINT,
    deleted_at TIMESTAMP
);

-- Auto-update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_tbl_users_updated_at
    BEFORE UPDATE ON tbl_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
