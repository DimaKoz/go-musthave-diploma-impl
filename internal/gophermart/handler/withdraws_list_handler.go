package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// WithdrawsListHandler handles GET `/api/user/withdrawals`.
func (h *BaseHandler) WithdrawsListHandler(ctx echo.Context) error {
	username := GetAuthFromCtx(ctx)
	withdraws, err := repository.FindWithdrawsByUsername(h.conn, username)
	if err != nil {
		var status int
		var logStr string
		if errors.Is(err, repository.ErrWithdrawsNoItems) {
			logStr = fmt.Sprintf("%s %s", "WithdrawsListHandler:", "no items:")
			status = http.StatusNoContent
		} else {
			logStr = fmt.Sprintf("%s %s %s", "WithdrawsListHandler:", "internal error:", err.Error())
			status = http.StatusInternalServerError
		}
		zap.S().Warn(logStr)
		_ = ctx.NoContent(status)

		return nil
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Response().WriteHeader(http.StatusOK)
	if err = json.NewEncoder(ctx.Response()).Encode(withdraws); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
