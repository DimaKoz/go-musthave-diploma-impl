package handler

import (
	"fmt"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/labstack/echo/v4"
)

// BaseHandler holds *pgx.Conn.
type BaseHandler struct {
	conn *sqldb.PgxIface
}

// NewBaseHandler returns a new BaseHandler.
func NewBaseHandler(dbConn *sqldb.PgxIface) *BaseHandler {
	return &BaseHandler{
		conn: dbConn,
	}
}

func AddAuthHeaders(ctx echo.Context, message string) {
	auth := fmt.Sprintf("Authorization:[%s]", message)
	ctx.Response().Header().Add("Authorization", auth)
	ctx.Response().Header().Add("Set-Cookie", auth)
}

func WrapHandlerErr(ctx echo.Context, statusCode int, msg string, errIn error) error {
	err := ctx.String(statusCode, fmt.Sprintf(msg, errIn))
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
