package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {
	if username == "" {
		return fmt.Errorf("no valid username")
	}
	cfg.CurrentUserName = username
	if err := write(*cfg); err != nil {
		return err
	}
	return nil
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

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}
