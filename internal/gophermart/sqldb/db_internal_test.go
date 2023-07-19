package sqldb

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
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

func TestFindOrderByNumberOk(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()

	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "username", "uploaded_at"}).
		AddRow("79927398713", "NEW", float32(0), "user1", now)

	mock.ExpectQuery(
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=\\$1").
		WithArgs("user1").
		WillReturnRows(rows)

	var pgConn PgxIface = mock
	cred, err := FindOrderByNumber(&pgConn, "user1")
	assert.NoError(t, err)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
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

const testUsernameDebit = "login2"

func TestGetDebitByUsernameReturnsErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var want float32
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(accrual\\),0\\)").
		WithArgs(testUsernameDebit).
		WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	cred, err := GetDebitByUsername(&pgConn, testUsernameDebit)
	assert.Error(t, err)
	assert.Equal(t, want, cred)
	assert.ErrorIs(t, err, io.EOF)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetDebitByUsernameReturns0(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var want float32
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(accrual\\),0\\)").
		WithArgs(testUsernameDebit).WillReturnRows(pgxmock.NewRows([]string{"accrual"}))

	var pgConn PgxIface = mock
	cred, err := GetDebitByUsername(&pgConn, testUsernameDebit)
	assert.NoError(t, err)
	assert.Equal(t, want, cred)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetDebitByUsernameReturns42(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	var want float32 = 42.0
	rows := mock.NewRows([]string{"accrual"}).
		AddRow(want)

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(accrual\\),0\\)").
		WithArgs(testUsernameDebit).WillReturnRows(rows)

	var pgConn PgxIface = mock
	cred, err := GetDebitByUsername(&pgConn, testUsernameDebit)
	assert.NoError(t, err)
	assert.Equal(t, want, cred)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetCreditByUsernameErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(testUsernameDebit).WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	_, err = GetCreditByUsername(&pgConn, testUsernameDebit)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetCreditByUsername0(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var want float32 = 42.0
	rows := mock.NewRows([]string{"sum"}).
		AddRow(want)

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(testUsernameDebit).WillReturnRows(rows)

	var pgConn PgxIface = mock
	cred, err := GetCreditByUsername(&pgConn, testUsernameDebit)
	assert.NoError(t, err)
	assert.Equal(t, want, cred)
	log.Println("cred:", cred)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetCreditByUsername42(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var want float32
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(testUsernameDebit).WillReturnRows(pgxmock.NewRows([]string{"sum"}))

	var pgConn PgxIface = mock
	cred, err := GetCreditByUsername(&pgConn, testUsernameDebit)
	assert.NoError(t, err)
	assert.Equal(t, want, cred)
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

func TestAddCredentialsErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectExec("insert into mart_users").
		WithArgs("user1", "pass1").
		WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	err = AddCredentials(&pgConn, credential.NewCredentials("user1", "pass1"))
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddOrder(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	mock.ExpectExec("insert into orders").
		WithArgs("79927398713", "NEW", float32(0), "user1", now).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	var pgConn PgxIface = mock
	order := accrual.NewOrderExt("79927398713", "NEW", float32(0), now, "user1")
	err = AddOrder(&pgConn, order)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddOrderErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	mock.ExpectExec("insert into orders").
		WithArgs("79927398713", "NEW", float32(0), "user1", now).WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	order := accrual.NewOrderExt("79927398713", "NEW", float32(0), now, "user1")
	err = AddOrder(&pgConn, order)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateOrder(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	mock.ExpectExec("UPDATE orders").
		WithArgs("NEW", float32(0), "79927398713").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	var pgConn PgxIface = mock
	order := accrual.NewOrderExt("79927398713", "NEW", float32(0), now, "user1")
	err = UpdateOrder(&pgConn, order)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateOrderErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	mock.ExpectExec("UPDATE orders").
		WithArgs("NEW", float32(0), "79927398713").
		WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	order := accrual.NewOrderExt("79927398713", "NEW", float32(0), now, "user1")
	err = UpdateOrder(&pgConn, order)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddWithdraw(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	mock.ExpectExec("insert into withdraws").
		WithArgs("79927398713", float32(0), "user1", now).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	var pgConn PgxIface = mock
	withdraw := accrual.NewWithdrawExt("79927398713", float32(0), now, "user1")
	err = AddWithdraw(&pgConn, *withdraw)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddWithdrawErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	mock.ExpectExec("insert into withdraws").
		WithArgs("79927398713", float32(0), "user1", now).
		WillReturnError(io.EOF)

	var pgConn PgxIface = mock
	withdraw := accrual.NewWithdrawExt("79927398713", float32(0), now, "user1")
	err = AddWithdraw(&pgConn, *withdraw)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
