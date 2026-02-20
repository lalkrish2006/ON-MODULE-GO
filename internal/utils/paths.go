package utils

import (
	"os"
	"path/filepath"
)

// ResolvePath attempts to find the absolute path for a given relative path
// by checking the current directory and progressively moving up to the project root.
func ResolvePath(relativePath string) string {
	// 1. Check if the path exists relative to the current working directory
	if _, err := os.Stat(relativePath); err == nil {
		return relativePath
	}

	// 2. Check assuming we are in cmd/server (typical mistake: go run main.go from subfolder)
	// Strategy: Walk up directories until we find "go.mod", which marks the root.
	
	dir, err := os.Getwd()
	if err != nil {
		return relativePath // Fallback
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found root
			return filepath.Join(dir, relativePath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root of filesystem without finding go.mod
			break
		}
		dir = parent
	}

	// Fallback: return original path if resolution fails
	return relativePath
}
