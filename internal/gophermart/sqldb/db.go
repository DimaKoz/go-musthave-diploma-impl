package sqldb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Close(context.Context) error
}

var errNoInfoConnectionDB = errors.New("no DB connection info")

// ConnectDB opens a connection to the database.
func ConnectDB(cfg *config.Config, logger echo.Logger) (*PgxIface, error) {
	inTestRunning := os.Getenv("GO_ENV1") == "testing"
	var conn PgxIface
	var err error
	if inTestRunning {
		conn, err = pgxmock.NewConn()
		if cfg == nil || cfg.ConnectionDB == "" {
			return nil, errNoInfoConnectionDB
		}
	} else {
		if cfg == nil || cfg.ConnectionDB == "" {
			return nil, errNoInfoConnectionDB
		}
		conn, err = pgx.Connect(context.Background(), cfg.ConnectionDB)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get a DB connection: %w", err)
	}
	timeout := 10

	if err = createTables(&conn, timeout); err != nil {
		return nil, err
	}
	logger.Info("successfully connected to db:", conn)

	return &conn, nil
}

func createTables(pgConn *PgxIface, timeout int) error {
	sqlString := `
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
CREATE INDEX IF NOT EXISTS idx_withdraws_status_username
    ON withdraws (username, sum);
COMMIT;
`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	if _, err := (*pgConn).Exec(ctx, sqlString); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func AddCredentials(pgConn *PgxIface, cred *credential.Credentials) error {
	_, err := (*pgConn).Exec(
		context.Background(),
		"insert into mart_users(name, password) values($1, $2)",
		cred.Username, cred.HashedPass)
	if err != nil {
		return fmt.Errorf("failed to insert into mart_users: %w", err)
	}

	return nil
}

func FindUserByUsername(pgConn *PgxIface, username string) (*credential.Credentials, error) {
	var cred *credential.Credentials
	var nameM, valueP string
	row := (*pgConn).QueryRow(context.Background(), "select name, password from mart_users where name=$1", username)
	err := row.Scan(&nameM, &valueP)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return cred, nil
		}

		return nil, fmt.Errorf("failed to scan a row: %w", err)
	}
	cred = credential.NewCredentials(nameM, valueP)

	return cred, nil
}

func AddOrder(pgConn *PgxIface, order *accrual.OrderExt) error {
	_, err := (*pgConn).Exec(
		context.Background(),
		"insert into orders(number, status, accrual, username, uploaded_at) values($1, $2, $3, $4, $5)",
		order.Number, order.Status, order.Accrual, order.Username, order.UploadedAt)
	if err != nil {
		return fmt.Errorf("failed to insert into orders: %w", err)
	}

	return nil
}

func AddWithdraw(pgConn *PgxIface, withdraw accrual.WithdrawExt) error {
	_, err := (*pgConn).Exec(
		context.Background(),
		"insert into withdraws(number, sum, username, processed_at) values($1, $2, $3, $4)",
		withdraw.Order, withdraw.Sum, withdraw.Username, withdraw.ProcessedAt)
	if err != nil {
		return fmt.Errorf("failed to insert into orders: %w", err)
	}

	return nil
}

func UpdateOrder(pgConn *PgxIface, order *accrual.OrderExt) error {
	_, err := (*pgConn).Exec(
		context.Background(),
		"UPDATE orders SET status = $1, accrual = $2 WHERE number = $3",
		order.Status, order.Accrual, order.Number)
	if err != nil {
		return fmt.Errorf("failed to update into orders: %w", err)
	}

	return nil
}

func FindOrderByNumber(pgConn *PgxIface, sNumber string) (*accrual.OrderExt, error) {
	var order *accrual.OrderExt
	var number, status, username string
	var uploadedAt time.Time
	var accrualV float32
	row := (*pgConn).QueryRow(context.Background(),
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=$1", sNumber)
	err := row.Scan(&number, &status, &accrualV, &username, &uploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order, nil
		}

		return nil, fmt.Errorf("failed to scan a row: %w", err)
	}
	order = accrual.NewOrderExt(number, status, accrualV, uploadedAt, username)

	return order, nil
}

func FindOrdersByUsername(pgConn *PgxIface, username string) (*[]accrual.OrderExt, error) {
	result := make([]accrual.OrderExt, 0)
	rows, err := (*pgConn).Query(context.Background(),
		"SELECT number, status, accrual, uploaded_at FROM orders WHERE username=$1", username)
	if err != nil {
		return &result, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var number, status string
		var uploadedAt time.Time
		var accrualV float32
		err = rows.Scan(&number, &status, &accrualV, &uploadedAt)
		if err != nil {
			return &result, fmt.Errorf("failed to scan a row: %w", err)
		}

		order := accrual.NewOrderExt(number, status, accrualV, uploadedAt, username)
		result = append(result, *order)
	}

	return &result, nil
}

func GetDebitByUsername(pgConn *PgxIface, username string) (float32, error) {
	row := (*pgConn).QueryRow(context.Background(),
		"SELECT COALESCE(SUM(accrual),0) FROM orders WHERE username=$1 AND status='PROCESSED'", username)
	var accrualV float32
	err := row.Scan(&accrualV)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to get debit: %w", err)
	}

	return accrualV, nil
}

func GetCreditByUsername(pgConn *PgxIface, username string) (float32, error) {
	row := (*pgConn).QueryRow(context.Background(),
		"SELECT COALESCE(SUM(sum),0) FROM withdraws WHERE username=$1", username)
	var accrualV float32
	err := row.Scan(&accrualV)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to get credit: %w", err)
	}

	return accrualV, nil
}

func FindWithdrawsByUsername(pgConn *PgxIface, username string) (*[]accrual.WithdrawExt, error) {
	result := make([]accrual.WithdrawExt, 0)
	rows, err := (*pgConn).Query(context.Background(),
		"SELECT number, sum, processed_at FROM withdraws WHERE username=$1 AND sum > 0.01 ORDER BY processed_at ASC",
		username)
	if err != nil {
		return &result, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var number string
		var processedAt time.Time
		var sum float32
		err = rows.Scan(&number, &sum, &processedAt)
		if err != nil {
			return &result, fmt.Errorf("failed to scan a row: %w", err)
		}

		withdraw := accrual.WithdrawExt{
			Order:       number,
			Sum:         sum,
			ProcessedAt: processedAt,
			Username:    username,
		}
		result = append(result, withdraw)
	}

	return &result, nil
}
