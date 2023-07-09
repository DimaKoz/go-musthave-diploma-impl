package util

import (
	"bytes"
	"os"

	"github.com/labstack/gommon/log"
)

// CaptureOutput catch logged info
// Use it with CaptureOutputCleanup.
func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()

	return buf.String()
}

// CaptureOutputCleanup restores logging into os.Stderr.
func CaptureOutputCleanup() {
	log.SetOutput(os.Stderr)
}
