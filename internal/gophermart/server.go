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
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/handler"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/middleware"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() {
	cfg := config.NewConfig()
	if err := setupConfig(cfg, config.ProcessEnvServer); err != nil {
		log.Error(err)

		return
	}
	echoFramework := echo.New()
	log.Info("cfg:" + cfg.String())
	var err error
	var conn *sqldb.PgxIface
	if conn, err = sqldb.ConnectDB(cfg, echoFramework.Logger); err == nil {
		defer (*conn).Close(context.Background())
	} else if os.Getenv("GO_ENV1") != "testing" {
		echoFramework.Logger.Errorf("failed to get a db connection by %s", err.Error())

		return
	}
	startServer(echoFramework, conn, *cfg)
}

var (
	errNoAddress = fmt.Errorf("server address is empty")
	errNoPathDB  = fmt.Errorf("db uri is empty")
)

func setupConfig(cfg *config.Config, processing config.ProcessEnv) error {
	if err := config.LoadConfig(cfg, processing); err != nil {
		return fmt.Errorf("couldn't create a config %w", err)
	}
	if cfg.Address == "" {
		return errNoAddress
	}
	if cfg.ConnectionDB == "" {
		return errNoPathDB
	}

	return nil
}

func startServer(echoFramework *echo.Echo, conn *sqldb.PgxIface, cfg config.Config) {
	// Setup
	baseHandler := handler.NewBaseHandler(conn, cfg)
	echoFramework.Logger.SetLevel(log.INFO)
	echoFramework.POST("/api/user/register", baseHandler.RegistrationHandler)
	echoFramework.POST("/api/user/login", baseHandler.LoginHandler)

	authM := middleware.AuthValidator(conn, echoFramework.Logger)
	orderValidM := middleware.OrderValidator(echoFramework.Logger)
	echoFramework.GET("/api/user/orders", baseHandler.OrdersListHandler, authM)
	echoFramework.POST("/api/user/orders", baseHandler.OrderUploadHandler, authM, orderValidM)

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
