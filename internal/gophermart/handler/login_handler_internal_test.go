package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testLoginRightBody = "{\"login\":\"login2\",\"password\":\"password2\"}"

func TestLoginHandler(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"}).
		AddRow("login2", "$2a$04$KujIDhc7zKDw0y2mVrNODOMYLBcc1B7kxTIiOf7unhaLHB/dr/9Mq")

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("login2").
		WillReturnRows(rs)

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

	baseH := NewBaseHandler(&pgConn, *cfg)

	err = baseH.LoginHandler(ctx)
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

func TestLoginHandlerBadRequest(t *testing.T) {
	// Mock echo
	echoFr := echo.New()
	body := "'"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := NewBaseHandler(nil, *cfg)

	err := baseH.LoginHandler(ctx)
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusBadRequest
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotHeaderA := got.Header.Get("Authorization")
	assert.Empty(t, gotHeaderA)
}

func TestLoginHandlerUnauthorizedUserNotFound(t *testing.T) {
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
		WithArgs("login5").
		WillReturnRows(rs)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := "{\"login\":\"login5\",\"password\":\"password5\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := NewBaseHandler(&pgConn, *cfg)

	err = baseH.LoginHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusUnauthorized
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotHeaderA := got.Header.Get("Authorization")
	assert.Empty(t, gotHeaderA)
}

func TestLoginHandlerUnauthorizedWrongPassword(t *testing.T) {
	// Mock db
	// DB connection
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, context.Background())

	rs := pgxmock.NewRows([]string{"name", "password"}).
		AddRow("login2", "wrong value")

	mock.ExpectQuery("select name, password from mart_users where name=\\$1").
		WithArgs("login2").
		WillReturnRows(rs)

	var pgConn sqldb.PgxIface = mock

	// Mock echo
	echoFr := echo.New()
	body := "{\"login\":\"login2\",\"password\":\"wrong value\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := NewBaseHandler(&pgConn, *cfg)

	err = baseH.LoginHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusUnauthorized
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotHeaderA := got.Header.Get("Authorization")
	assert.Empty(t, gotHeaderA)
}

func TestLoginHandlerInternalErr(t *testing.T) {
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
	body := "{\"login\":\"login2\",\"password\":\"wrong value\"}"
	req := httptest.NewRequest(echo.POST, "http://localhost:1323/admin/user_points/settings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := echoFr.NewContext(req, rec)
	defer echoFr.Close()

	cfg := config.NewConfig()

	baseH := NewBaseHandler(&pgConn, *cfg)

	err = baseH.LoginHandler(ctx)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	got := rec.Result()
	defer got.Body.Close()

	wantStatusCode := http.StatusInternalServerError
	assert.Equal(t, wantStatusCode, got.StatusCode, "StatusCode got: %v, want: %v", got.StatusCode, wantStatusCode)

	gotHeaderA := got.Header.Get("Authorization")
	assert.Empty(t, gotHeaderA)
}
