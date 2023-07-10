package sqldb

import (
	"context"
	"fmt"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDBConnectGetLogger(t *testing.T) *echo.Logger {
	t.Helper()

	logger := echo.New().Logger

	return &logger
}

func TestConnectDBErrNoConnection1(t *testing.T) {
	logger := testDBConnectGetLogger(t)
	cfg := config.NewConfig()
	cfg.ConnectionDB = "***"
	conn, err := ConnectDB(cfg, *logger)
	assert.Nil(t, conn)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid dsn")
}

func TestConnectDBErrNoConnection(t *testing.T) {
	logger := testDBConnectGetLogger(t)

	conn, err := ConnectDB(nil, *logger)
	assert.Nil(t, conn)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNoInfoConnectionDB)
}

func TestCreateTables(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	var pgConn PgxIface = mock
	timeout := 10

	result := pgconn.NewCommandTag("CREATE TABLE")
	mock.
		ExpectExec("CREATE TABLE IF NOT EXISTS mart_users").
		WillReturnResult(result)

	err = createTables(&pgConn, timeout)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		assert.Error(t, err, fmt.Sprintf("there were unfulfilled expectations: %s", err))
	}
}

func TestCreateTablesErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	var pgConn PgxIface = mock
	timeout := 10

	err = createTables(&pgConn, timeout)
	assert.Error(t, err)
}
