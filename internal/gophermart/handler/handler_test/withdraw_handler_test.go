package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const loginNameTestingWithdraw = "login2"

func TestWithdrawHandler(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(accrual\\),0\\)").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(mock.NewRows([]string{"accrual"}).AddRow(float32(44)))

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(pgxmock.NewRows([]string{"sum"}).AddRow(float32(2)))

	mock.ExpectExec("insert into withdraws").
		WithArgs("2377225624", float32(2), loginNameTestingWithdraw, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	bodyStr := "{\"order\": \"2377225624\",\n    \"sum\": 2\n}"
	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", strings.NewReader(bodyStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func TestWithdrawHandlerBadRequest(t *testing.T) { // Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	bodyStr := "'"
	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", strings.NewReader(bodyStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(nil, *cfg)

	err := baseH.WithdrawHandler(ctx)
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusUnprocessableEntity
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func TestWithdrawHandlerNoMoney(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(accrual\\),0\\)").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(mock.NewRows([]string{"accrual"}).AddRow(float32(44)))

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(pgxmock.NewRows([]string{"sum"}).AddRow(float32(2)))

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	bodyStr := "{\"order\": \"2377225624\",\n    \"sum\": 9999\n}"
	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", strings.NewReader(bodyStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusPaymentRequired
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func TestWithdrawHandlerInternalErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(accrual\\),0\\)").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(mock.NewRows([]string{"accrual"}).AddRow(float32(44)))

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs(loginNameTestingWithdraw).
		WillReturnRows(pgxmock.NewRows([]string{"sum"}).AddRow(float32(2)))

	mock.ExpectExec("insert into withdraws").
		WithArgs("2377225624", float32(2), loginNameTestingWithdraw, pgxmock.AnyArg()).
		WillReturnError(io.EOF)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	bodyStr := "{\"order\": \"2377225624\",\n    \"sum\": 2\n}"
	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", strings.NewReader(bodyStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", loginNameTestingWithdraw))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.WithdrawHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}
