CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS tbl_users (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    user_name TEXT NOT NULL,
    email TEXT,
    password TEXT NOT NULL,
    role_name TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    login_session TEXT,
    last_login TIMESTAMPTZ,
    currency_id INTEGER,
    language_id INTEGER,
    status_id INTEGER,
    "order" INTEGER,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT,
    updated_at TIMESTAMPTZ,
    deleted_by BIGINT,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tbl_users_user_name
    ON tbl_users (user_name);

CREATE OR REPLACE FUNCTION set_tbl_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_tbl_users_updated_at ON tbl_users;
CREATE TRIGGER trg_tbl_users_updated_at
    BEFORE UPDATE ON tbl_users
    FOR EACH ROW
    EXECUTE FUNCTION set_tbl_users_updated_at();

CREATE TABLE IF NOT EXISTS auth_users (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    refresh_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ
);

INSERT INTO tbl_users (
    first_name,
    last_name,
    user_name,
    email,
    password,
    role_name,
    role_id,
    is_admin,
    created_by
)
VALUES (
    'System',
    'Admin',
    'ADMIN',
    'admin@example.com',
    crypt('admin123', gen_salt('bf')),
    'ADMIN',
    1,
    TRUE,
    NULL
)
ON CONFLICT (user_name) DO NOTHING;
