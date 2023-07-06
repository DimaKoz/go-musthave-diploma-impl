package main

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGracefulShutdown(t *testing.T) {
	go func() { // killer
		time.Sleep(5 * time.Second)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT) // syscall.SIGTERM
		assert.NoError(t, err)
	}()
	main()
}
