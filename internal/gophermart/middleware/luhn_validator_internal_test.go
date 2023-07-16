package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestOrderValidatorMiddleware(t *testing.T) {
	echoFramework := echo.New()
	echoFramework.Use(OrderValidator(echoFramework.Logger))
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	// Using the ServerHTTP on echo will trigger the router and middleware
	echoFramework.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t, rec.Body.String(), "{\"message\":\"bad order\"}\n")
}

const (
	badOrder = "5262148203"
	okOrder  = "5262148207"
)

func TestOrderValidatorMiddleware404(t *testing.T) {
	echoFramework := echo.New()
	echoFramework.Use(OrderValidator(echoFramework.Logger))
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(okOrder))
	rec := httptest.NewRecorder()

	// Using the ServerHTTP on echo will trigger the router and middleware
	echoFramework.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, rec.Body.String(), "{\"message\":\"Not Found\"}\n")
}

func TestOrderValidatorMiddleware422(t *testing.T) {
	echoFramework := echo.New()
	echoFramework.Use(OrderValidator(echoFramework.Logger))
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(badOrder))
	rec := httptest.NewRecorder()

	// Using the ServerHTTP on echo will trigger the router and middleware
	echoFramework.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t, rec.Body.String(), "{\"message\":\"bad order\"}\n")
}
