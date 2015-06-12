package util

import (
	"os"
	"path/filepath"
)

var cutoroot string = filepath.Join(os.Getenv("GOPATH"), "cutoroot")

const osname = "linux"

func init() {
	os.Setenv("CUTOROOT", cutoroot)
}
