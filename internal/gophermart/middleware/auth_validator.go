package middleware

import (
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AuthValidator checks 'Authorization' header and its value.
func AuthValidator(dbConn *sqldb.PgxIface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			return validate(echoCtx, dbConn, next)
		}
	}
}

func validate(echoCtx echo.Context, dbConn *sqldb.PgxIface, next echo.HandlerFunc) error {
	zap.S().Info("AuthValidator: ", echoCtx.Request().Method, " ", echoCtx.Request().URL)
	zap.S().Info("AuthValidator: Headers: ", echoCtx.Request().Header)
	isAuthorized := handler.IsAuthorized(echoCtx, dbConn)
	authHeader := echoCtx.Request().Header.Get("Authorization")
	if !isAuthorized {
		zap.S().Warn("AuthValidator: failed to check an authorization for: [%s]",
			authHeader)
		err := handler.WrapHandlerErr(echoCtx, http.StatusUnauthorized,
			"AuthValidator: failed to check an authorization: %s", handler.ErrUnauthorised)

		return fmt.Errorf("%w", err)
	}

	zap.S().Infof("AuthValidator: Authorization header is correct for [%s]",
		authHeader)
	if err := next(echoCtx); err != nil {
		echoCtx.Error(err)
	}

	return nil
}
