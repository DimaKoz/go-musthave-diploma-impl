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
	middleware2 "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func Run() {
	var err error
	loggerZap := zap.Must(zap.NewDevelopment())

	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(loggerZap)

	zap.ReplaceGlobals(loggerZap)

	cfg := config.NewConfig()
	if err = setupConfig(cfg, config.ProcessEnvServer); err != nil {
		log.Error(err)
		zap.S().Error(err)

		return
	}
	echoFramework := echo.New()
	zap.S().Info("cfg:" + cfg.String())

	var conn *sqldb.PgxIface
	if conn, err = sqldb.ConnectDB(cfg, echoFramework.Logger); err == nil {
		defer (*conn).Close(context.Background())
	} else if os.Getenv("GO_ENV1") != "testing" {
		zap.S().Errorf("failed to get a db connection by %s", err.Error())

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
	loggerConfig := middleware.GetRequestLoggerConfig()
	log2 := middleware2.RequestLoggerWithConfig(loggerConfig)
	log3 := middleware2.BodyDump(middleware.GetBodyLoggerHandler())

	// Setup
	baseHandler := handler.NewBaseHandler(conn, cfg)
	echoFramework.Logger.SetLevel(log.INFO)
	echoFramework.POST("/api/user/register", baseHandler.RegistrationHandler,
		log2, log3)
	echoFramework.POST("/api/user/login", baseHandler.LoginHandler,
		log2, log3)

	authM := middleware.AuthValidator(conn)

	echoFramework.GET("/api/user/orders", baseHandler.OrdersListHandler,
		log2, log3, authM)
	echoFramework.POST("/api/user/orders", baseHandler.OrderUploadHandler,
		log2, log3, authM, middleware.OrderValidator())
	echoFramework.POST("/api/user/balance/withdraw", baseHandler.WithdrawHandler,
		log2, log3, authM)
	echoFramework.GET("/api/user/balance", baseHandler.BalanceHandler,
		log2, log3, authM)
	echoFramework.GET("/api/user/withdrawals", baseHandler.WithdrawsListHandler,
		log2, log3, authM)

	// Start server
	go func(cfg config.Config) {
		zap.S().Info("start server")
		if err := echoFramework.Start(cfg.Address); err != nil && errors.Is(err, http.ErrServerClosed) {
			echoFramework.Logger.Warn("shutting down the server")
		}
	}(cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	zap.S().Info("quit...")
	timeoutDelay := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutDelay)*time.Second)
	defer cancel()
	if err := echoFramework.Shutdown(ctx); err != nil {
		zap.S().Fatal(err)
	}
}
