package gophermart

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() {
	cfg := config.NewConfig()
	err := setupConfig(cfg, config.ProcessEnvServer)
	if err != nil {
		log.Error(err)

		return
	}
	log.Info("cfg:" + cfg.String())
	startServer(*cfg)
}

var errNoAddress = fmt.Errorf("server address is empty")

func setupConfig(cfg *config.Config, processing config.ProcessEnv) error {
	if err := config.LoadConfig(cfg, processing); err != nil {
		return fmt.Errorf("couldn't create a config %w", err)
	}
	if cfg.Address == "" {
		return errNoAddress
	}

	return nil
}

func startServer(cfg config.Config) {
	// Setup
	echoFramework := echo.New()
	echoFramework.Logger.SetLevel(log.INFO)
	echoFramework.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK) //nolint:wrapcheck
	})
	// Start server
	go func(cfg config.Config) {
		echoFramework.Logger.Info("start server")
		if err := echoFramework.Start(cfg.Address); err != nil && errors.Is(err, http.ErrServerClosed) {
			echoFramework.Logger.Warn("shutting down the server")
		}
	}(cfg)

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
