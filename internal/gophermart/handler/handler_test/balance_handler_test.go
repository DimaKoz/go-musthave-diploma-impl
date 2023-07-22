package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceHandler(t *testing.T) {
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
		WithArgs("login2").
		WillReturnRows(mock.NewRows([]string{"accrual"}).AddRow(float32(44)))

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs("login2").
		WillReturnRows(pgxmock.NewRows([]string{"sum"}).AddRow(float32(2)))

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.BalanceHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	want := "{\"current\":42,\"withdrawn\":2}\n"
	gotBody, err := io.ReadAll(got.Body)
	require.NoError(t, err)

	assert.Equal(t, want, string(gotBody), "Body got: %v, want: %v", string(gotBody), want)
}

func TestBalanceHandlerInternalErr(t *testing.T) {
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
		WithArgs("login2").
		WillReturnError(io.EOF)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.BalanceHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
}

func TestBalanceHandlerErr(t *testing.T) {
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
		WithArgs("login2").
		WillReturnRows(mock.NewRows([]string{"accrual"}).AddRow(float32(44)))

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(sum\\),0\\)").
		WithArgs("login2").
		WillReturnRows(pgxmock.NewRows([]string{"sum"}).AddRow(float32(2)))

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	ctx.Response().Writer = &testresponsewriter{w: rec} //nolint:exhaustruct

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.BalanceHandler(ctx)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
