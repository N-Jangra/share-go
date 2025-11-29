package utils

import (
	"os"
)

// EnsureDir creates the directory if it doesn't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
