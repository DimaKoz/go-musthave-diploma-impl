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

/*const (
	alreadyUploadedByOwner   = http.StatusOK                  // 200 — номер заказа уже был загружен этим пользователем.
	orderAccepted            = http.StatusAccepted            // 202 — новый номер заказа принят в обработку
	badRequest               = http.StatusBadRequest          // 400 — неверный формат запроса
	alreadyUploadedByAnother = http.StatusConflict            // 409 — номер заказа уже был загружен другим пользователем
	internalError            = http.StatusInternalServerError // 500 — внутренняя ошибка сервера
)
*/

// OrderUploadHandler handles POST `/api/user/orders`.
func (h *BaseHandler) OrderUploadHandler(ctx echo.Context) error {
	reqBody := []byte{}
	if ctx.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(ctx.Request().Body)
	}
	log.Println("OrderUploadHandler:", "body:", string(reqBody))
	var acc accrual.OrderAccrual
	httpc := resty.New().SetBaseURL(h.cfg.Accrual)

	req := httpc.R().
		SetResult(&acc).
		SetPathParam("number", string(reqBody))

	resp, err := req.Post("/api/orders/{number}")
	if resp != nil && resp.Request != nil {
		log.Println("OrderUploadHandler:", "req.URL:", resp.Request.URL)
	}

	if err != nil {
		log.Println("OrderUploadHandler:", "req.Get err:", err)
	}
	if resp != nil {
		log.Println("OrderUploadHandler:", "req.Get StatusCode:", resp.StatusCode())
		log.Println("OrderUploadHandler:", "req.Get resp:", resp.String())
	}
	log.Println("OrderUploadHandler:", "acc:", acc)
	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
