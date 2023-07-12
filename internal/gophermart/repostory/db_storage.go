package repostory

import (
	"fmt"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/sqldb"
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
