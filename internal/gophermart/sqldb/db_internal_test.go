package sqldb

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
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

func TestFindUserByUsernameReturnsUser(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"}).AddRow("user1", "pass1")

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("user1").
		WillReturnRows(rs)

	var pgConn PgxIface = mock
	cred, err := FindUserByUsername(&pgConn, "user1")
	assert.NoError(t, err)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindUserByUsernameReturnsNil(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"})

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("user2").
		WillReturnRows(rs)

	var pgConn PgxIface = mock
	cred, err := FindUserByUsername(&pgConn, "user2")
	assert.NoError(t, err)
	assert.Nil(t, cred)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindUserByUsernameReturnsErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("user2").
		WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	cred, err := FindUserByUsername(&pgConn, "user2")
	assert.Error(t, err)
	assert.Nil(t, cred)
	assert.ErrorIs(t, err, io.EOF)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddCredentials(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectExec("insert into mart_users").
		WithArgs("user1", "pass1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	var pgConn PgxIface = mock
	err = AddCredentials(&pgConn, credential.NewCredentials("user1", "pass1"))
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
