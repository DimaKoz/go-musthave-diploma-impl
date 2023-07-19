BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS mart_users
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(72) NOT NULL,
    password VARCHAR(72) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mart_users_name
    ON mart_users USING hash (name);

CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    number      VARCHAR(42)      NOT NULL,
    status      VARCHAR(10)      NOT NULL,
    accrual     DOUBLE PRECISION NOT NULL,
    username    VARCHAR(72)      NOT NULL,
    uploaded_at TIMESTAMP        NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_number
    ON orders USING hash (number);

CREATE INDEX IF NOT EXISTS idx_orders_username
    ON orders USING hash (username);

CREATE INDEX IF NOT EXISTS idx_orders_status
    ON orders USING hash (status);

CREATE INDEX IF NOT EXISTS idx_orders_status_username
    ON orders (username, status);

CREATE TABLE IF NOT EXISTS withdraws
(
    id           SERIAL PRIMARY KEY,
    number       VARCHAR(42)      NOT NULL,
    sum          DOUBLE PRECISION NOT NULL,
    username     VARCHAR(72)      NOT NULL,
    processed_at TIMESTAMP        NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_withdraws_number
    ON withdraws USING hash (number);

CREATE INDEX IF NOT EXISTS idx_withdraws_username
    ON withdraws USING hash (username);

COMMIT;