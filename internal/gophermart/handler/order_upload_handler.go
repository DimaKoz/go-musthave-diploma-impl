package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// OrderUploadHandler handles POST `/api/user/orders`.
func (h *BaseHandler) OrderUploadHandler(ctx echo.Context) error {
	reqBody := []byte{}
	if ctx.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(ctx.Request().Body)
	}
	log.Println("OrderUploadHandler:", "body:", string(reqBody))

	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
