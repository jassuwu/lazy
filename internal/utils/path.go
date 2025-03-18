package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandTilde replaces the ~ character at the beginning of a path with the user's home directory
// Automatically handles all the OSs
func ExpandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path // Return original path if we can't get home dir
	}

	return filepath.Join(home, path[1:])
}
