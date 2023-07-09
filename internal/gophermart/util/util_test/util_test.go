package util_test

import (
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/util"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	t.Cleanup(func() {
		util.CaptureOutputCleanup()
	})
	want := "abc"

	got := util.CaptureOutput(func() {
		log.Info("abc")
	})
	assert.Contains(t, got, want)
}
