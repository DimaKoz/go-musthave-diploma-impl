package sqldb

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConnectDBErrNoConnection1(t *testing.T) {
	cfg := config.NewConfig()
	cfg.ConnectionDB = "***"
	conn, err := ConnectDB(cfg)
	assert.Nil(t, conn)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid dsn")
}

func TestConnectDBErrNoConnection(t *testing.T) {
	conn, err := ConnectDB(nil)
	assert.Nil(t, conn)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNoInfoConnectionDB)
}

func TestConnectDBErrNoConnection2(t *testing.T) {
	_ = os.Setenv("GO_ENV1", "testing") //nolint:tenv
	defer os.Unsetenv("GO_ENV1")

	conn, err := ConnectDB(nil)
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
		WithArgs("79927398713").
		WillReturnRows(rows)

	var pgConn PgxIface = mock
	order, err := FindOrderByNumber(&pgConn, "79927398713")
	assert.NoError(t, err)
	zap.S().Infoln("order:", order)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	rows2 := pgxmock.NewRows([]string{"number", "status", "accrual", "uploaded_at"}).
		AddRow("79927398713", "NEW", float32(0), now)

	mock.ExpectQuery(
		"SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("user1").
		WillReturnRows(rows2)
	orders, err := FindOrdersByUsername(&pgConn, "user1")
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Len(t, *orders, 1)
	zap.S().Infoln("orders:", orders)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindOrderByNumberNoRowsErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "username", "uploaded_at"})

	mock.ExpectQuery(
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=\\$1").
		WithArgs("79927398713").
		WillReturnRows(rows)

	var pgConn PgxIface = mock
	order, err := FindOrderByNumber(&pgConn, "79927398713")
	assert.NoError(t, err)
	assert.Nil(t, order)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindOrderByNumberScanErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "username", "uploaded_at"}).
		AddRow("79927398713", "NEW", float32(0), "user1", pgxmock.AnyArg())

	mock.ExpectQuery(
		"SELECT number, status, accrual, username, uploaded_at FROM orders WHERE number=\\$1").
		WithArgs("79927398713").
		WillReturnRows(rows)

	var pgConn PgxIface = mock
	order, err := FindOrderByNumber(&pgConn, "79927398713")
	assert.Error(t, err)
	assert.Nil(t, order)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindOrdersByUsernameScanErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var pgConn PgxIface = mock

	rows2 := pgxmock.NewRows([]string{"number", "status", "accrual", "uploaded_at"}).
		AddRow("79927398713", "NEW", float32(0), pgxmock.AnyArg())

	mock.ExpectQuery(
		"SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("user1").
		WillReturnRows(rows2)
	orders, err := FindOrdersByUsername(&pgConn, "user1")
	assert.Error(t, err)
	assert.NotNil(t, orders)
	assert.Len(t, *orders, 0)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindOrdersByUsernameErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var pgConn PgxIface = mock

	mock.ExpectQuery(
		"SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("user1").
		WillReturnError(io.EOF)
	orders, err := FindOrdersByUsername(&pgConn, "user1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
	assert.NotNil(t, orders)
	assert.Len(t, *orders, 0)
	zap.S().Infoln("orders:", orders)

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
	zap.S().Infoln("cred:", cred)

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
	zap.S().Infoln("cred:", cred)

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
	zap.S().Infoln("cred:", cred)

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
	zap.S().Infoln("cred:", cred)

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
	zap.S().Infoln("cred:", cred)

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
	zap.S().Infoln("cred:", cred)

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

func TestGetCreditByUsername42(t *testing.T) {
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
	zap.S().Infoln("cred:", cred)

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

	var want float32
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(testUsernameDebit).WillReturnRows(pgxmock.NewRows([]string{"sum"}))

	var pgConn PgxIface = mock
	cred, err := GetCreditByUsername(&pgConn, testUsernameDebit)
	assert.NoError(t, err)
	assert.Equal(t, want, cred)
	zap.S().Infoln("cred:", cred)

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

func TestFindWithdrawsByUsernameOk(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()

	var pgConn PgxIface = mock

	rows2 := pgxmock.NewRows([]string{"number", "sum", "processed_at"}).
		AddRow("79927398713", float32(0), now)

	mock.ExpectQuery(
		"SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs("user1").
		WillReturnRows(rows2)
	withdraws, err := FindWithdrawsByUsername(&pgConn, "user1")
	assert.NoError(t, err)
	assert.NotNil(t, withdraws)
	assert.Len(t, *withdraws, 1)
	zap.S().Infoln("withdraws:", withdraws)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindWithdrawsByUsernameErrScan(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var pgConn PgxIface = mock

	rows2 := pgxmock.NewRows([]string{"number", "sum", "processed_at"}).
		AddRow("79927398713", float32(0), pgxmock.AnyArg())

	mock.ExpectQuery(
		"SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs("user1").
		WillReturnRows(rows2)
	withdraws, err := FindWithdrawsByUsername(&pgConn, "user1")
	assert.Error(t, err)
	assert.NotNil(t, withdraws)
	assert.Len(t, *withdraws, 0)
	zap.S().Infoln("withdraws:", withdraws)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindWithdrawsByUsernameErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))

	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	var pgConn PgxIface = mock

	mock.ExpectQuery(
		"SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs("user1").
		WillReturnError(io.EOF)
	withdraws, err := FindWithdrawsByUsername(&pgConn, "user1")
	assert.Error(t, err)
	assert.NotNil(t, withdraws)
	assert.Len(t, *withdraws, 0)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
