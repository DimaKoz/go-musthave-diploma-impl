package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// OrdersListHandler handles GET `/api/user/orders`.
func (h *BaseHandler) OrdersListHandler(ctx echo.Context) error {
	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
