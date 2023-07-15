BEGIN TRANSACTION;

DROP INDEX IF EXISTS idx_mart_users_name;

DROP INDEX IF EXISTS idx_orders_number;

DROP INDEX IF EXISTS idx_orders_username;

DROP INDEX IF EXISTS idx_orders_status;

DROP TABLE IF EXISTS mart_users;

DROP TABLE IF EXISTS orders;

COMMIT;