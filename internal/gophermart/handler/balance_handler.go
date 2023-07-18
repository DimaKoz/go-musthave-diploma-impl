package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// BalanceHandler handles POST `/api/user/balance`.
func (h *BaseHandler) BalanceHandler(ctx echo.Context) error {
	username := GetAuthFromCtx(ctx)
	log.Println("BalanceHandler:", "username:", username)

	balance, err := repository.GetBalance(h.conn, username)
	if err != nil {
		zap.S().Warnf("BalanceHandler: internal error %s", err.Error())
		if err = ctx.NoContent(http.StatusInternalServerError); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	if err = ctx.JSON(http.StatusOK, balance); err != nil {
		err = fmt.Errorf("%w", err)

		return err
	}

	return nil
}
