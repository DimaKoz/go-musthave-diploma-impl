package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
)

// OrdersListHandler handles GET `/api/user/orders`.
func (h *BaseHandler) OrdersListHandler(ctx echo.Context) error {
	username := GetAuthFromCtx(ctx)
	orders, err := repository.GetOrdersByUser(h.conn, username)
	if err != nil {
		logStr := fmt.Sprintf("%s %s %s", "OrdersListHandler:", "internal error:", err.Error())
		log.Println(logStr)
		if err := ctx.String(http.StatusInternalServerError, logStr); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	if orders == nil || len(*orders) == 0 {
		logStr := fmt.Sprintf("%s %s", "OrdersListHandler:", "no data to response")
		log.Println(logStr)
		if err := ctx.String(http.StatusNoContent, logStr); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Response().WriteHeader(http.StatusOK)
	if err := json.NewEncoder(ctx.Response()).Encode(orders); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
