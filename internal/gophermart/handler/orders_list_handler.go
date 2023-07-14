package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

// OrdersListHandler handles GET `/api/user/orders`.
func (h *BaseHandler) OrdersListHandler(ctx echo.Context) error {
	reqBody := []byte{}
	if ctx.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(ctx.Request().Body)
	}
	log.Println("OrdersListHandler:", "RequestURI:", ctx.Request().URL)
	log.Println("OrdersListHandler:", "Header:", ctx.Request().Header)
	log.Println("OrdersListHandler:", "body:", string(reqBody))
	ctx.Request().Body.Close()

	var acc accrual.Order
	httpc := resty.New().SetHostURL(h.cfg.Accrual) //nolint:staticcheck

	req := httpc.R().
		SetResult(&acc).
		SetPathParam("number", string(reqBody))

	resp, err := req.Get("/api/orders/{number}")
	if resp != nil && resp.Request != nil {
		log.Println("OrdersListHandler:", "req.URL:", resp.Request.URL)
	}

	if err != nil {
		log.Println("OrdersListHandler:", "req.Get err:", err)
	}
	if resp != nil {
		log.Println("OrdersListHandler:", "req.Get StatusCode:", resp.StatusCode())
		log.Println("OrdersListHandler:", "req.Get resp:", resp.String())
	}
	log.Println("OrdersListHandler:", "acc:", acc)
	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
