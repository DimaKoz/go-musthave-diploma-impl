package sqldb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
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
func ConnectDB(cfg *config.Config, logger echo.Logger) (*pgx.Conn, error) {
	if cfg == nil || cfg.ConnectionDB == "" {
		return nil, errNoInfoConnectionDB
	}
	conn, err := pgx.Connect(context.Background(), cfg.ConnectionDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get a DB connection: %w", err)
	}
	timeout := 10
	var db PgxIface = conn
	if err = createTables(&db, timeout); err != nil {
		return nil, err
	}
	logger.Info("successfully connected to db:", conn)

	return conn, nil
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

COMMIT;`

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
