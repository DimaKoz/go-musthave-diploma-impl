package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCredentialsReturnsErrUserNameNotFound(t *testing.T) {
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

	var pgConn sqldb.PgxIface = mock

	cred, err := GetCredentials(&pgConn, "user2")

	assert.ErrorIs(t, err, ErrUserNameNotFound)
	assert.Nil(t, cred)
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

	var pgConn sqldb.PgxIface = mock
	err = AddCredentials(&pgConn, *credential.NewCredentials("user1", "pass1"))
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddNewOrderErr(t *testing.T) {
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

	var pgConn sqldb.PgxIface = mock
	err = AddNewOrder(&pgConn, "79927398713", "user1")
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestWithdrawErr(t *testing.T) {
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
		WithArgs("login2").WillReturnRows(rows)
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs("login2").WillReturnError(io.EOF)

	var pgConn sqldb.PgxIface = mock
	withdrawInternal := accrual.WithdrawAccrual{} //nolint:exhaustruct
	withdraw := withdrawInternal.GetWithdrawExt("login2", time.Now())
	err = ProcessWithdraw(&pgConn, *withdraw)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
