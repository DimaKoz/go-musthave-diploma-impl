package security

import "github.com/ShiraazMoollatjie/goluhn"

func IsValidLuhnNumber(luhnNumber string) bool {
	return goluhn.Validate(luhnNumber) == nil
}
