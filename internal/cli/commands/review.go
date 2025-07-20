package commands

import (
	"fmt"
	"os"
)

func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}
	return string(data), nil
}
