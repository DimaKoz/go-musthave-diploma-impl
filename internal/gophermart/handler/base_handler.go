package handler

import (
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
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
