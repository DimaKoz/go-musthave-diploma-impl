package middleware

import (
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
)

// AuthValidator checks 'Authorization' header and its value.
func AuthValidator(dbConn *sqldb.PgxIface, logger echo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			isAuthorized := handler.IsAuthorized(echoCtx, dbConn)
			authHeader := echoCtx.Request().Header.Get("Authorization")
			if !isAuthorized {
				logger.Warnf("AuthValidator: failed to check an authorization for: [%s]",
					authHeader)
				err := handler.WrapHandlerErr(echoCtx, http.StatusUnauthorized,
					"AuthValidator: failed to check an authorization: %s", handler.ErrUnauthorised)

				return fmt.Errorf("%w", err)
			}

			logger.Infof("AuthValidator: Authorization header is correct for [%s]",
				authHeader)
			if err := next(echoCtx); err != nil {
				echoCtx.Error(err)
			}

			return nil
		}
	}
}
