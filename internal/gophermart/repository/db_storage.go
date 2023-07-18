package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/accrual"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
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

var ErrWithdrawsNoItems = fmt.Errorf("there are not withdrawals")

func FindWithdrawsByUsername(pgConn *sqldb.PgxIface, username string) (*[]accrual.WithdrawExt, error) {
	withdraws, err := sqldb.FindWithdrawsByUsername(pgConn, username)
	if err != nil {
		err = fmt.Errorf("failed to get withdraws by: %w", err)

		return nil, err
	}
	if len(*withdraws) == 0 {
		return nil, ErrWithdrawsNoItems
	}

	return withdraws, nil
}

var (
	ErrCantAddOrder                = fmt.Errorf("failed to add order")
	ErrOrderAlreadyExistsByOwner   = fmt.Errorf("failed to add order: order already exists by the owner")
	ErrOrderAlreadyExistsByAnother = fmt.Errorf("failed to add order: order already exists by another user")

	ErrWithdrawNoMoney = fmt.Errorf("failed to process withdrawal: no money")
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

func GetOrdersByUser(pgConn *sqldb.PgxIface, username string) (*[]accrual.OrderExt, error) {
	orders, err := sqldb.FindOrdersByUsername(pgConn, username)
	if err != nil {
		return orders, fmt.Errorf("%w", err)
	}

	return orders, nil
}

func GetBalance(pgConn *sqldb.PgxIface, username string) (accrual.BalanceExt, error) {
	result := accrual.BalanceExt{Current: 0, Withdrawn: 0}
	debit, err := sqldb.GetDebitByUsername(pgConn, username)
	zap.S().Debugln("debit:", debit, "err:", err)
	if err != nil {
		return result, fmt.Errorf("%w", err)
	}
	credit, err := sqldb.GetCreditByUsername(pgConn, username)
	zap.S().Debugln("credit:", credit, "err:", err)
	if err != nil {
		return result, fmt.Errorf("%w", err)
	}
	result.Current = debit - credit
	result.Withdrawn = credit

	return result, nil
}

func ProcessWithdraw(pgConn *sqldb.PgxIface, withdraw accrual.WithdrawExt) error {
	balance, err := GetBalance(pgConn, withdraw.Username)
	if err != nil {
		return fmt.Errorf("withdraw: failed to get balance order by:%w", err)
	}
	if balance.Current <= 0 || withdraw.Sum > balance.Current {
		return ErrWithdrawNoMoney
	}

	if err = sqldb.AddWithdraw(pgConn, withdraw); err != nil {
		return fmt.Errorf("withdraw: failed to add withdraw by:%w", err)
	}

	return nil
}
