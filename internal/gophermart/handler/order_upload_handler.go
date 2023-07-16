package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repostory"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

const (
	alreadyUploadedByOwner = http.StatusOK       // 200 — номер заказа уже был загружен этим пользователем.
	orderAccepted          = http.StatusAccepted // 202 — новый номер заказа принят в обработку
	// badRequest               = http.StatusBadRequest          // 400 — неверный формат запроса.
	alreadyUploadedByAnother = http.StatusConflict            // 409 — номер заказа уже был загружен другим пользователем
	internalError            = http.StatusInternalServerError // 500 — внутренняя ошибка сервера
)

// OrderUploadHandler handles POST `/api/user/orders`.
func (h *BaseHandler) OrderUploadHandler(ctx echo.Context) error {
	reqBody := []byte{}
	if ctx.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(ctx.Request().Body)
	}
	orderNumber := string(reqBody)
	log.Println("OrderUploadHandler:", "body:", orderNumber)
	username := GetAuthFromCtx(ctx)
	err := repostory.AddNewOrder(h.conn, orderNumber, username)

	var respStatus int

	switch {
	case err == nil:
		respStatus = orderAccepted
	case errors.Is(err, repostory.ErrOrderAlreadyExistsByOwner):
		respStatus = alreadyUploadedByOwner
	case errors.Is(err, repostory.ErrOrderAlreadyExistsByAnother):
		respStatus = alreadyUploadedByAnother
	default:
		respStatus = internalError
	}

	/*	if err == nil {
			respStatus = orderAccepted
		} else if errors.Is(err, repostory.ErrOrderAlreadyExistsByOwner) {
			respStatus = alreadyUploadedByOwner
		} else if errors.Is(err, repostory.ErrOrderAlreadyExistsByAnother) {
			respStatus = alreadyUploadedByAnother
		} else {
			respStatus = internalError
		}
	*/
	go sendAccRequest(orderNumber, h.cfg.Accrual)
	if err := ctx.NoContent(respStatus); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func sendAccRequest(number string, baseURL string) {
	var acc accrual.OrderAccrual
	httpc := resty.New().SetBaseURL(baseURL)

	req := httpc.R().
		SetResult(&acc).
		SetPathParam("number", number)

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
}
