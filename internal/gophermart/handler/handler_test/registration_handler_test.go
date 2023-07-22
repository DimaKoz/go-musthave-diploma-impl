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

func TestRegistrationHandler(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"})

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("login2").
		WillReturnRows(rs)

	mock.ExpectExec("insert into mart_users").
		WithArgs("login2", pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := testLoginRightBody
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.RegistrationHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusOK
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotHeaderA := got.Header.Get("Authorization")
	wantHeaderA := "Authorization:[login2]"
	assert.Equal(t, wantHeaderA, gotHeaderA)
}

func TestRegistrationHandlerAddCredentialsErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"})

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("login2").
		WillReturnRows(rs)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := "{\"login\":\"login2\",\"password\":\"password2\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := echoFr.NewContext(req, rec)

	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.RegistrationHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	echoFr.ServeHTTP(rec, req)

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, rec.Code)
	res := rec.Result()
	gotHeaderA := res.Header.Get("Authorization")
	err = res.Body.Close()
	require.NoError(t, err)
	assert.Empty(t, gotHeaderA)
}

func TestRegistrationHandlerBadCredentialsErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"})

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("login2").
		WillReturnRows(rs)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := "{\"login\":\"login2\"," +
		"\"password\":\"password2password2password2password2password2password2password2password2password2\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := echoFr.NewContext(req, rec)

	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.RegistrationHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	echoFr.ServeHTTP(rec, req)

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, rec.Code)
	res := rec.Result()
	gotHeaderA := res.Header.Get("Authorization")
	err = res.Body.Close()
	require.NoError(t, err)
	assert.Empty(t, gotHeaderA)
}

func TestRegistrationHandlerNoUserErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rows := pgxmock.
		NewRows([]string{"name", "password"}).
		AddRow("login2", "password2")

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := "{\"login\":\"login2\"," +
		"\"password\":\"password2\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := echoFr.NewContext(req, rec)

	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.RegistrationHandler(ctx)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	echoFr.ServeHTTP(rec, req)

	wantStatusCode := http.StatusConflict
	assert.Equal(t, wantStatusCode, rec.Code)
	res := rec.Result()
	gotHeaderA := res.Header.Get("Authorization")
	err = res.Body.Close()
	require.NoError(t, err)
	assert.Empty(t, gotHeaderA)
}

func TestRegistrationHandlerBadRequest(t *testing.T) {
	// Mock echo
	echoFr := echo.New()
	body := "'"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := echoFr.NewContext(req, rec)

	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(nil, *cfg)

	err := baseH.RegistrationHandler(ctx)
	assert.NoError(t, err)

	wantStatusCode := http.StatusBadRequest
	assert.Equal(t, wantStatusCode, rec.Code)
	res := rec.Result()
	gotHeaderA := res.Header.Get("Authorization")
	err = res.Body.Close()
	require.NoError(t, err)
	assert.Empty(t, gotHeaderA)
}

func TestRegistrationHandlerUnknownCredentialsErr(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("login2").
		WillReturnError(io.EOF)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := "{\"login\":\"login2\"," +
		"\"password\":\"password2\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := echoFr.NewContext(req, rec)

	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := handler.NewBaseHandler(&pgConn, *cfg)

	err = baseH.RegistrationHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	echoFr.ServeHTTP(rec, req)

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, rec.Code)
	res := rec.Result()
	gotHeaderA := res.Header.Get("Authorization")
	err = res.Body.Close()
	require.NoError(t, err)
	assert.Empty(t, gotHeaderA)
}
