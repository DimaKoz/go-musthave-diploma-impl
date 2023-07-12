package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/repostory"
	"github.com/labstack/echo/v4"
)

// RegistrationHandler handles `/api/user/register`.
func (h *BaseHandler) RegistrationHandler(ctx echo.Context) error {
	incomeCred := &credential.IncomeCredentials{} //nolint:exhaustruct
	if err := ctx.Bind(incomeCred); err != nil {
		return wrapRegHandlerErr(ctx, http.StatusBadRequest, "RegistrationHandler: failed to parse json: %s", err)
	}

	cred, err := repostory.GetCredentials(h.conn, incomeCred.Login)
	if cred != nil {
		return fmt.Errorf("%w", ctx.String(http.StatusConflict,
			fmt.Sprintf("RegistrationHandler:  [%s] login is already taken", incomeCred.Login)))
	}
	if err != nil && !errors.Is(err, repostory.ErrUserNameNotFound) {
		return wrapRegHandlerErr(ctx, http.StatusInternalServerError, "RegistrationHandler: failed to find login by: %s", err)
	}

	cred = credential.NewCredentials(incomeCred.Login, "")
	if err = cred.HashPass(incomeCred.Password); err != nil {
		return wrapRegHandlerErr(ctx, http.StatusInternalServerError,
			"RegistrationHandler: failed to hash the pass by: %s", err)
	}

	if err = repostory.AddCredentials(h.conn, *cred); err != nil {
		return wrapRegHandlerErr(ctx, http.StatusInternalServerError,
			"RegistrationHandler: failed to save the credentials by: %s", err)
	}

	auth := fmt.Sprintf("Authorization:[%s]", incomeCred.Login)
	ctx.Response().Header().Add("Authorization", auth)
	ctx.Response().Header().Add("Set-Cookie", auth)
	if err := ctx.NoContent(http.StatusOK); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func wrapRegHandlerErr(ctx echo.Context, statusCode int, msg string, errIn error) error {
	err := ctx.String(statusCode, fmt.Sprintf(msg, errIn))
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
