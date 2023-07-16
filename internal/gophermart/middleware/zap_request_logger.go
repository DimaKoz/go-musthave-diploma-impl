package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var zapSugar zap.SugaredLogger

var logValuesFunc = func(c echo.Context, loggerValues middleware.RequestLoggerValues) error {
	zapSugar.Infow("request",
		zap.String("Method", loggerValues.Method),
		zap.String("URI", loggerValues.URI),
		zap.Duration("latency", loggerValues.Latency),
	)
	zapSugar.Infow("response",
		zap.Int("status", loggerValues.Status),
		zap.String("length", loggerValues.ContentLength),
		zap.Int64("size", loggerValues.ResponseSize),
	)

	return nil
}

// GetRequestLoggerConfig returns middleware.RequestLoggerConfig.
func GetRequestLoggerConfig(sugar zap.SugaredLogger) middleware.RequestLoggerConfig {
	zapSugar = sugar
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
func GetBodyLoggerHandler(sugar zap.SugaredLogger) middleware.BodyDumpHandler {
	return func(c echo.Context, reqBody, resBody []byte) {
		sugar.Infow(
			"body:", "reqBody:", string(reqBody),
		)
	}
}
