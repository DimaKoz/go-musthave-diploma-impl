package credential_test

import (
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/model/credential"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewCredentials(t *testing.T) {
	username := "a"
	pass := "b"
	want := &credential.Credentials{Username: username, HashedPass: pass}
	got := credential.NewCredentials(username, pass)
	assert.Equal(t, want, got)
}

func TestCredentials(t *testing.T) {
	passBefore := "message1"
	username := "user1"
	cred := credential.NewCredentials(username, "")
	err := cred.HashPass(passBefore)
	assert.NoError(t, err)
	assert.True(t, cred.IsPassCorrect(passBefore))
}

func TestCredentials1(t *testing.T) {
	tooLongPassBefore := "message1message1message1message1message1message1message" +
		"1message1message1message1message1message1message1message1message1message1message1" +
		"message1message1message1message1message1message1message1message1message1"
	username := "user1"
	cred := credential.NewCredentials(username, "")
	err := cred.HashPass(tooLongPassBefore)
	assert.Error(t, err)
	assert.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
}
