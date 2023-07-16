package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLogValuesFunc(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	sugar := *logger.Sugar()
	zapSugar = sugar
	e := echo.New()
	assert.NoError(t, logValuesFunc(e.AcquireContext(), middleware.RequestLoggerValues{})) //nolint:exhaustruct
}

func TestGetRequestLoggerConfig(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	sugar := *logger.Sugar()

	type args struct {
		sugar zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want middleware.RequestLoggerConfig
	}{
		{
			name: "test RequestLoggerConfig ",
			args: args{
				sugar: sugar,
			},
			want: middleware.RequestLoggerConfig{ //nolint:exhaustruct
				LogURI:           true,
				LogStatus:        true,
				LogLatency:       true,
				LogContentLength: true,
				LogResponseSize:  true,
				LogMethod:        true,
				LogValuesFunc:    logValuesFunc,
			},
		},
	}

	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			got := GetRequestLoggerConfig(test.args.sugar)
			assert.Equal(t, got.LogMethod, test.want.LogMethod)
			assert.Equal(t, got.LogURI, test.want.LogURI)
			assert.Equal(t, got.LogStatus, test.want.LogStatus)
			assert.Equal(t, got.LogResponseSize, test.want.LogResponseSize)
			assert.Equal(t, got.LogLatency, test.want.LogLatency)
		})
	}
}
