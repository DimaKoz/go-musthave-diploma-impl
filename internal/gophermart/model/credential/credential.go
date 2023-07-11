package credential

import (
	"fmt"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/security"
)

type Credentials struct {
	Username   string
	HashedPass string
}

func NewCredentials(username string, hashedPassword string) *Credentials {
	return &Credentials{Username: username, HashedPass: hashedPassword}
}

func (cred *Credentials) HashPass(password string) error {
	var err error
	var hash string
	if hash, err = security.GetHash(password); err == nil {
		cred.HashedPass = hash
	} else {
		return fmt.Errorf("failed to hash the password, error: %w", err)
	}

	return nil
}

func (cred *Credentials) IsPassCorrect(password string) bool {
	return security.IsRightHash(password, cred.HashedPass)
}
