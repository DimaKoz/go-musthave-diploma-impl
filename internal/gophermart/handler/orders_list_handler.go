package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// OrdersListHandler handles GET `/api/user/orders`.
func (h *BaseHandler) OrdersListHandler(ctx echo.Context) error {
	username := GetAuthFromCtx(ctx)
	orders, err := repository.GetOrdersByUser(h.conn, username)
	if err != nil {
		logStr := fmt.Sprintf("%s %s %s", "OrdersListHandler:", "internal error:", err.Error())
		zap.S().Infoln(logStr)
		_ = ctx.NoContent(http.StatusInternalServerError)

		return nil
	}

	if orders == nil || len(*orders) == 0 {
		logStr := fmt.Sprintf("%s %s", "OrdersListHandler:", "no data to response")
		zap.S().Infoln(logStr)
		_ = ctx.NoContent(http.StatusNoContent)

		return nil
	}
	checkOrders(h, orders)

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Response().WriteHeader(http.StatusOK)
	if err = json.NewEncoder(ctx.Response()).Encode(orders); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func checkOrders(baseH *BaseHandler, orders *[]accrual.OrderExt) {
	for idx, order := range *orders {
		if order.Status != accrual.OrderStatusInvalid &&
			order.Status != accrual.OrderStatusProcessed {
			continue
		}
		newOrder := SendAccRequest(baseH.conn, order.Number, baseH.cfg.Accrual, order.Username)
		if newOrder != nil {
			(*orders)[idx].Accrual = newOrder.Accrual
			(*orders)[idx].Status = newOrder.Status
		}
	}
}
