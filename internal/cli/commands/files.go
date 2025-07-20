package commands

import (
	"fmt"
	"os"
)

// ListDirectoryContents reads the contents of a single directory.
func ListDirectoryContents(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}
	return files, nil
}
