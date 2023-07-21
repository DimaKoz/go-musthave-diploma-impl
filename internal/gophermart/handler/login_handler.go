package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
)

// LoginHandler handles `/api/user/login`.
func (h *BaseHandler) LoginHandler(ctx echo.Context) error {
	incomeCred := &credential.IncomeCredentials{} //nolint:exhaustruct
	if err := ctx.Bind(incomeCred); err != nil {
		return WrapHandlerErr(ctx, http.StatusBadRequest, "LoginHandler: failed to parse json: %s", err)
	}

	cred, err := repository.GetCredentials(h.conn, incomeCred.Login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNameNotFound) {
			return WrapHandlerErr(ctx, http.StatusUnauthorized,
				"LoginHandler: failed to find login by: %s", err)
		}

		return WrapHandlerErr(ctx, http.StatusInternalServerError,
			"LoginHandler: failed to find login by: %s", err)
	}

	if !cred.IsPassCorrect(incomeCred.Password) {
		errStr := fmt.Sprintf(" [%s/%s] wrong credentials",
			incomeCred.Login, incomeCred.Password)
		err = errors.New(errStr) //nolint:goerr113

		return WrapHandlerErr(ctx, http.StatusUnauthorized,
			"LoginHandler: failed to login by: %s", err)
	}

	AddAuthHeaders(ctx, incomeCred.Login)
	_ = ctx.NoContent(http.StatusOK)

	return nil
}
