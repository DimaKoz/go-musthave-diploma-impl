package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	withdrawAccepted      = http.StatusOK                  // 200 — успешная обработка запроса.
	withdrawNotEnough     = http.StatusPaymentRequired     // 402 — на счету недостаточно средств
	withdrawInternalError = http.StatusInternalServerError // 500 — внутренняя ошибка сервера
)

// WithdrawHandler handles POST `/api/user/balance/withdraw`.
func (h *BaseHandler) WithdrawHandler(ctx echo.Context) error {
	withdrawInternal := accrual.WithdrawAccrual{} //nolint:exhaustruct
	if err := ctx.Bind(&withdrawInternal); err != nil {
		return WrapHandlerErr(ctx, http.StatusUnprocessableEntity,
			"WithdrawalHandler: failed to get withdrawInternal: %s", fmt.Errorf("%w", err))
	}
	zap.S().Infoln("WithdrawalHandler:", "withdrawInternal:", withdrawInternal)
	username := GetAuthFromCtx(ctx)

	withdraw := withdrawInternal.GetWithdrawExt(username, time.Now())

	err := repository.ProcessWithdraw(h.conn, *withdraw)

	var respStatus int

	switch {
	case err == nil:
		respStatus = withdrawAccepted
	case errors.Is(err, repository.ErrWithdrawNoMoney):
		respStatus = withdrawNotEnough
	default:
		respStatus = withdrawInternalError
	}

	_ = ctx.NoContent(respStatus)

	return nil
}
