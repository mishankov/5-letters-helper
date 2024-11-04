package testing

import (
	"bytes"
	"log"
	"os"
)

func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
