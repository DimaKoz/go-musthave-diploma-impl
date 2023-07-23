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

func TestOrdersListHandler(t *testing.T) {
	// Mock DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "uploaded_at"}).
		AddRow("79927398713", "NEW", float32(0), now).
		AddRow("79927398714", "PROCESSED", float32(42), now)

	mock.ExpectQuery("SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("login2").
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrdersListHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)
	marshaledNow, err := now.MarshalText()
	require.NoError(t, err)
	order1 := fmt.Sprintf("{\"number\":\"79927398713\",\"status\":\"NEW\",\"accrual\":0,\"uploaded_at\":\"%s\"}",
		string(marshaledNow))
	order2 := fmt.Sprintf("{\"number\":\"79927398714\",\"status\":\"PROCESSED\",\"accrual\":42,\"uploaded_at\":\"%s\"}",
		string(marshaledNow))
	want := fmt.Sprintf("[%s,%s]\n", order1, order2)
	gotBody, err := io.ReadAll(got.Body)
	require.NoError(t, err)
	assert.Equal(t, want, string(gotBody), "Body got: %v, want: %v", string(gotBody), want)
}

func TestOrdersListHandlerStatus500(t *testing.T) {
	// Mock DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("login2").
		WillReturnError(http.ErrAbortHandler)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrdersListHandler(ctx)
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

func TestOrdersListHandlerStatusNoContent(t *testing.T) {
	// Mock DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "uploaded_at"})

	mock.ExpectQuery("SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("login2").
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrdersListHandler(ctx)
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

func TestOrdersListHandlerResponseErr(t *testing.T) {
	// Mock DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())
	now := time.Now()
	rows := pgxmock.NewRows([]string{"number", "status", "accrual", "uploaded_at"}).
		AddRow("79927398713", "NEW", float32(0), now).
		AddRow("79927398714", "PROCESSED", float32(42), now)

	mock.ExpectQuery("SELECT number, status, accrual, uploaded_at FROM orders WHERE username=\\$1").
		WithArgs("login2").
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	defer echoFr.Close()

	req := httptest.NewRequest(echo.GET, "http://localhost:1323/", nil)
	req.Header.Add("Authorization",
		fmt.Sprintf("Authorization:[%s]", "login2"))

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	ctx.Response().Writer = &testresponsewriter{w: rec} //nolint:exhaustruct

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.OrdersListHandler(ctx)
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
