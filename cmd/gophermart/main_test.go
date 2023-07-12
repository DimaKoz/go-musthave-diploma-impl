package main

import (
	"io"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGracefulShutdown(t *testing.T) {
	rescueStdout := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer
	_ = os.Setenv("GO_ENV1", "testing") //nolint:tenv
	defer os.Unsetenv("GO_ENV1")
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
		util.CaptureOutputCleanup()
	})
	var wGroup sync.WaitGroup
	wGroup.Add(1)

	go func() {
		defer wGroup.Done()
		main()
		time.Sleep(10 * time.Second)
	}()
	wGroup.Wait()

	writer.Close()
	out, _ := io.ReadAll(reader)
	reader.Close()
	os.Stdout = rescueStdout

	assert.Contains(t, string(out), "shutting down the server")
}
