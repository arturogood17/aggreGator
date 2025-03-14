package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) ([]byte, error) {
	cfg.CurrentUserName = username
	jsonFile, err := json.Marshal(cfg)
	if err != nil {
		return []byte{}, err
	}
	return jsonFile, nil
}

func Read() (Config, error) {
	file, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
