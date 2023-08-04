package handler

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository/cooldown"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
	var reqBody []byte
	if ctx.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(ctx.Request().Body)
	}
	orderNumber := string(reqBody)
	zap.S().Infoln("OrderUploadHandler:", "body:", orderNumber)
	username := GetAuthFromCtx(ctx)
	err := repository.AddNewOrder(h.conn, orderNumber, username)

	var respStatus int

	switch {
	case err == nil:
		respStatus = orderAccepted
	case errors.Is(err, repository.ErrOrderAlreadyExistsByOwner):
		respStatus = alreadyUploadedByOwner
	case errors.Is(err, repository.ErrOrderAlreadyExistsByAnother):
		respStatus = alreadyUploadedByAnother
	default:
		respStatus = internalError
	}

	go SendAccRequest(h.conn, orderNumber, h.cfg.Accrual, username)
	_ = ctx.NoContent(respStatus)

	return nil
}

func sleepIfCooldown() {
	for { // we are waiting for finished cooldown
		if cooldown.IsAccrualReady() {
			break
		}
		time.Sleep(time.Second)
	}
}

func SendAccRequest(pgConn *sqldb.PgxIface, number string, baseURL string, username string) *accrual.OrderExt {
	var acc accrual.OrderAccrual
	logger := zap.S()
	httpc := resty.New().SetBaseURL(baseURL)

	sleepIfCooldown()

	for {
		req := httpc.R().
			SetResult(&acc).
			SetPathParam("number", number)

		resp, err := req.Get("/api/orders/{number}")
		if resp != nil && resp.Request != nil {
			logger.Info("OrderUploadHandler:", "req.URL:", resp.Request.URL)
			logger.Info("OrderUploadHandler:", "req.Get StatusCode:", resp.StatusCode())
			logger.Info("OrderUploadHandler:", "req.Get resp:", resp.String())
		}

		if err != nil {
			logger.Info("OrderUploadHandler:", "req.Get err:", err)
		}

		if resp.StatusCode() == http.StatusTooManyRequests {
			cooldown.NeedAccrualCooldown()
			time.Sleep(1 * time.Minute)

			continue
		}

		logger.Infoln("OrderUploadHandler:", "acc:", acc)
		if acc.Order != "" {
			order := acc.GetOrderExt(username, time.Now())
			err = sqldb.UpdateOrder(pgConn, order)
			if err != nil {
				logger.Warn("err:", err.Error())
			} else {
				return order
			}
		}

		break
	}

	return nil
}
