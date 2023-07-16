package repostory

import (
	"errors"
	"fmt"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/jackc/pgx/v5"
)

func AddCredentials(pgConn *sqldb.PgxIface, cred credential.Credentials) error {
	err := sqldb.AddCredentials(pgConn, &cred)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to add credentials by: %w", err)
}

var ErrUserNameNotFound = fmt.Errorf("username not found")

func GetCredentials(pgConn *sqldb.PgxIface, username string) (*credential.Credentials, error) {
	cred, err := sqldb.FindUserByUsername(pgConn, username)
	if cred == nil && err == nil {
		err = ErrUserNameNotFound
	}

	return cred, err
}

var (
	ErrCantAddOrder                = fmt.Errorf("failed to add order")
	ErrOrderAlreadyExistsByOwner   = fmt.Errorf("failed to add order: order already exists by the owner")
	ErrOrderAlreadyExistsByAnother = fmt.Errorf("failed to add order: order already exists by another user")
)

func AddNewOrder(pgConn *sqldb.PgxIface, sNumber string, username string) error {
	order, err := sqldb.FindOrderByNumber(pgConn, sNumber)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return ErrCantAddOrder
	}

	if order != nil {
		if order.Username == username {
			return ErrOrderAlreadyExistsByOwner
		}

		return ErrOrderAlreadyExistsByAnother
	}

	order = accrual.NewOrderExt(sNumber, accrual.OrderStatusNew, 0, time.Now(), username)

	if err = sqldb.AddOrder(pgConn, order); err != nil {
		return fmt.Errorf("failed to add order by:%w", err)
	}

	return nil
}
