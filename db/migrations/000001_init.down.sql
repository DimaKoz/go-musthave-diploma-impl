BEGIN TRANSACTION;

DROP INDEX IF EXISTS idx_mart_users_name;

DROP TABLE IF EXISTS mart_users;

COMMIT;