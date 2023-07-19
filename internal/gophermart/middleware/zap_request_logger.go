package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var logValuesFunc = func(c echo.Context, loggerValues middleware.RequestLoggerValues) error {
	zap.S().Infow("request",
		zap.String("Method", loggerValues.Method),
		zap.String("URI", loggerValues.URI),
		zap.Duration("latency", loggerValues.Latency),
	)
	zap.S().Infow("response",
		zap.Int("status", loggerValues.Status),
		zap.String("length", loggerValues.ContentLength),
		zap.Int64("size", loggerValues.ResponseSize),
	)

	return nil
}

// GetRequestLoggerConfig returns middleware.RequestLoggerConfig.
func GetRequestLoggerConfig() middleware.RequestLoggerConfig {
	result := middleware.RequestLoggerConfig{ //nolint:exhaustruct
		LogURI:           true,
		LogStatus:        true,
		LogLatency:       true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogMethod:        true,
		LogValuesFunc:    logValuesFunc,
		LogHeaders:       []string{},
	}

	return result
}

// GetBodyLoggerHandler returns middleware.BodyDumpHandler.
func GetBodyLoggerHandler() middleware.BodyDumpHandler {
	return func(c echo.Context, reqBody, resBody []byte) {
		zap.S().Infow(
			"body:", "reqBody:", string(reqBody),
		)
	}
}
