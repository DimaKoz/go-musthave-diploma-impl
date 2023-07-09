package main

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGracefulShutdown(t *testing.T) {
	go func() { // killer
		time.Sleep(5 * time.Second)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT) // syscall.SIGTERM
		assert.NoError(t, err)
	}()
	origValue := os.Getenv(config.EnvKeyAddress)
	err := os.Setenv(config.EnvKeyAddress, ":8080") //nolint:tenv
	require.NoError(t, err)
	origValue1 := os.Getenv(config.EnvKeyDatabaseURI)
	err = os.Setenv(config.EnvKeyDatabaseURI, ":8080") //nolint:tenv
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Setenv(config.EnvKeyAddress, origValue)
		_ = os.Setenv(config.EnvKeyDatabaseURI, origValue1)
	})
	output := util.CaptureOutput(func() {
		main()
	})
	assert.Contains(t, output, "shutting down the server")
}
