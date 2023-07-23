package handler

import (
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// BalanceHandler handles GET `/api/user/balance`.
func (h *BaseHandler) BalanceHandler(ctx echo.Context) error {
	username := GetAuthFromCtx(ctx)
	zap.S().Infoln("BalanceHandler:", "username:", username)

	balance, err := repository.GetBalance(h.conn, username)
	if err != nil {
		zap.S().Warnf("BalanceHandler: internal error %s", err.Error())
		_ = ctx.NoContent(http.StatusInternalServerError)

		return nil
	}

	if err = ctx.JSON(http.StatusOK, balance); err != nil {
		err = fmt.Errorf("%w", err)

		return err
	}

	return nil
}
