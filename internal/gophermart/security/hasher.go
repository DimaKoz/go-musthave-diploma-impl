package security

import "golang.org/x/crypto/bcrypt"

func GetHash(message string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(message), bcrypt.MinCost)

	return string(bytes), err
}

func IsRightHash(message, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(message))

	return err == nil
}
