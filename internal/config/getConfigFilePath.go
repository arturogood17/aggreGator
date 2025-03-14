package config

import (
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	file, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(file, fileName), nil
}
