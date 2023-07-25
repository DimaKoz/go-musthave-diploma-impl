package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawsListHandler(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	rows := pgxmock.NewRows([]string{"number", "sum", "processed_at"}).
		AddRow("79927398713", float32(0), now)

	mock.ExpectQuery("SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawsListHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
	marshaledNow, err := now.MarshalText()
	require.NoError(t, err)
	want := "[{\"order\":\"79927398713\",\"sum\":0,\"processed_at\":\"" + string(marshaledNow) + "\"}]\n"
	gotBody, err := io.ReadAll(got.Body)
	require.NoError(t, err)
	assert.Equal(t, want, string(gotBody), "Body got: %v, want: %v", string(gotBody), want)
}

func TestWithdrawsListHandlerResponseErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	rows := pgxmock.NewRows([]string{"number", "sum", "processed_at"}).
		AddRow("79927398713", float32(0), now)

	mock.ExpectQuery("SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	ctx.Response().Writer = &testresponsewriter{w: rec} //nolint:exhaustruct
	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawsListHandler(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
	gotBody, err := io.ReadAll(got.Body)
	require.NoError(t, err)
	assert.Empty(t, gotBody)
}

func TestWithdrawsListHandlerEmptyResult(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rows := pgxmock.NewRows([]string{"number", "sum", "processed_at"})

	mock.ExpectQuery("SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawsListHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusNoContent
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotBody, err := io.ReadAll(got.Body)
	require.NoError(t, err)
	assert.Empty(t, gotBody)
}

func TestWithdrawsListHandlerInternalErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT number, sum, processed_at FROM withdraws WHERE username=\\$1").
		WithArgs(loginNameTestingWithdraw).
		WillReturnError(http.ErrAbortHandler)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawsListHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotBody, err := io.ReadAll(got.Body)
	require.NoError(t, err)
	assert.Empty(t, gotBody)
}
