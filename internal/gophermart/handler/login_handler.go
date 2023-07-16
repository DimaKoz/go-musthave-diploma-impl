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

	if cred == nil {
		return fmt.Errorf("%w", ctx.String(http.StatusUnauthorized,
			fmt.Sprintf("LoginHandler: [%s] failed to find login", incomeCred.Login)))
	}

	if !cred.IsPassCorrect(incomeCred.Password) {
		return fmt.Errorf("%w", ctx.String(http.StatusUnauthorized,
			fmt.Sprintf("LoginHandler: [%s/%s] wrong credentials", incomeCred.Login, incomeCred.Password)))
	}

	AddAuthHeaders(ctx, incomeCred.Login)
	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
