package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/security"
	"github.com/labstack/echo/v4"
)

// OrderValidator checks luhn number for order value.
func OrderValidator(logger echo.Logger) echo.MiddlewareFunc {
	badLuhn := echo.NewHTTPError(http.StatusUnprocessableEntity, "bad order")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			// Request
			reqBody := []byte{}
			if echoCtx.Request().Body != nil { // Read
				reqBody, _ = io.ReadAll(echoCtx.Request().Body)
			}
			echoCtx.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset
			if len(reqBody) == 0 {
				logger.Warnf("failed to find order number")

				return badLuhn
			}
			isValid := security.IsValidLuhnNumber(string(reqBody))
			if !isValid {
				logger.Warnf("bad lunh number for [%s]", string(reqBody))

				return badLuhn
			}
			logger.Infof("lunh number is correct for [%s]", string(reqBody))
			if err := next(echoCtx); err != nil {
				echoCtx.Error(err)
			}

			return nil
		}
	}
}
