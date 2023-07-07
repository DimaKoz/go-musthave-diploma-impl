package security_test

import (
	"log"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/security"
	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	message := "message1"
	gotHash, err := security.GetHash(message)
	got := security.IsRightHash(message, gotHash)
	log.Println("message:[" + message + "], hash:[" + gotHash + "]")
	assert.NoError(t, err, "GetHash() error = %v, wantErr %v", err, nil)
	assert.True(t, got, "GetHash() got = %v", gotHash)
}
