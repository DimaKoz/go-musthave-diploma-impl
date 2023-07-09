package util

import (
	"bytes"
	"os"

	"github.com/labstack/gommon/log"
)

func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)

	return buf.String()
}
