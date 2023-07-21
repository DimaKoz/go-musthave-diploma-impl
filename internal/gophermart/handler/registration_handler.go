package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repository"
	"github.com/labstack/echo/v4"
)

// RegistrationHandler handles `/api/user/register`.
func (h *BaseHandler) RegistrationHandler(ctx echo.Context) error {
	incomeCred := &credential.IncomeCredentials{} //nolint:exhaustruct
	if err := ctx.Bind(incomeCred); err != nil {
		return WrapHandlerErr(ctx, http.StatusBadRequest, "RegistrationHandler: failed to parse json: %s", err)
	}

	cred, err := repository.GetCredentials(h.conn, incomeCred.Login)
	if cred != nil {
		return fmt.Errorf("%w", ctx.String(http.StatusConflict,
			fmt.Sprintf("RegistrationHandler:  [%s] login is already taken", incomeCred.Login)))
	}
	if err != nil && !errors.Is(err, repository.ErrUserNameNotFound) {
		return WrapHandlerErr(ctx, http.StatusInternalServerError, "RegistrationHandler: failed to find login by: %s", err)
	}

	cred = credential.NewCredentials(incomeCred.Login, "")
	if err = cred.HashPass(incomeCred.Password); err != nil {
		return WrapHandlerErr(ctx, http.StatusInternalServerError,
			"RegistrationHandler: failed to hash the pass by: %s", err)
	}

	if err = repository.AddCredentials(h.conn, *cred); err != nil {
		return WrapHandlerErr(ctx, http.StatusInternalServerError,
			"RegistrationHandler: failed to save the credentials by: %s", err)
	}

	AddAuthHeaders(ctx, incomeCred.Login)
	_ = ctx.NoContent(http.StatusOK)

	return nil
}
