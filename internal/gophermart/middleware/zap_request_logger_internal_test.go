package middleware

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestLogValuesFunc(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
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
			name: "testRequestLoggerConfig",
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
			got := GetRequestLoggerConfig()
			assert.Equal(t, got.LogMethod, test.want.LogMethod)
			assert.Equal(t, got.LogURI, test.want.LogURI)
			assert.Equal(t, got.LogStatus, test.want.LogStatus)
			assert.Equal(t, got.LogResponseSize, test.want.LogResponseSize)
			assert.Equal(t, got.LogLatency, test.want.LogLatency)
		})
	}
}

func TestGetBodyLoggerHandler(t *testing.T) {
	want := "body:[TestBody]"

	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(func(e zapcore.Entry) error {
		assert.Equal(t, want, e.Message)

		return nil
	})))
	original := zap.L()
	zap.ReplaceGlobals(logger)
	t.Cleanup(func() {
		zap.ReplaceGlobals(original)
	})
	var err error
	echoFramework := echo.New()
	defer func(echoFr *echo.Echo) {
		err = echoFr.Close()
		require.NoError(t, err)
	}(echoFramework)

	echoFramework.Use(middleware.BodyDump(GetBodyLoggerHandler()))

	req := httptest.NewRequest(echo.GET, "/", strings.NewReader("TestBody"))
	rec := httptest.NewRecorder()

	// Using the ServerHTTP on echo will trigger the router and middleware
	echoFramework.ServeHTTP(rec, req)
}
