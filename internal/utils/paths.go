package utils

import (
	"os"
	"path/filepath"
)

// Not strictly needed now, but could hold path-related utilities.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "bootstrapme"), nil
}
