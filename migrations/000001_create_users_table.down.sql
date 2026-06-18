-- 000001_create_users_table.down.sql
-- Reverts the users table creation

DROP TRIGGER IF EXISTS trg_tbl_users_updated_at ON tbl_users;
DROP FUNCTION IF EXISTS update_updated_at();
DROP TABLE IF EXISTS tbl_users;
