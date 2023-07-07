package gophermart

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() {
	startServer()
}

func startServer() {
	// Setup
	echoFramework := echo.New()
	echoFramework.Logger.SetLevel(log.INFO)
	echoFramework.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK) //nolint:wrapcheck
	})
	// Start server
	go func() {
		echoFramework.Logger.Info("start server")
		if err := echoFramework.Start(":8080"); err != nil && errors.Is(err, http.ErrServerClosed) {
			echoFramework.Logger.Warn("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	echoFramework.Logger.Info("quit...")
	timeoutDelay := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutDelay)*time.Second)
	defer cancel()
	if err := echoFramework.Shutdown(ctx); err != nil {
		echoFramework.Logger.Fatal(err)
	}
}
