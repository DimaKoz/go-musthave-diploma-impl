package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthValidatorOk(t *testing.T) {
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

	echoFramework := echo.New()
	defer func(echoFr *echo.Echo) {
		err = echoFr.Close()
		require.NoError(t, err)
	}(echoFramework)

	echoFramework.Use(AuthValidator(&pgConn))
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Authorization:[%s]", "login2"))
	rec := httptest.NewRecorder()

	// Using the ServerHTTP on echo will trigger the router and middleware
	echoFramework.ServeHTTP(rec, req)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAuthValidator401(t *testing.T) {
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

	echoFramework := echo.New()
	defer func(echoFr *echo.Echo) {
		err = echoFr.Close()
		require.NoError(t, err)
	}(echoFramework)

	echoFramework.Use(AuthValidator(&pgConn))
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Authorization:[%s]", "login2"))
	rec := httptest.NewRecorder()

	// Using the ServerHTTP on echo will trigger the router and middleware
	echoFramework.ServeHTTP(rec, req)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, rec.Body.String(), "AuthValidator: failed to check an authorization: unauthorized request")
}
