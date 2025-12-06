package runtimeutil

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

// FindProjectRoot walks up the directory tree to find the project root (where go.mod exists).
func FindProjectRoot() string {
	dir := lo.Must(os.Getwd())

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod not found)")
		}
		dir = parent
	}
}
