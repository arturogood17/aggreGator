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

func Read() (Config, error) {
	file, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	file += "/.gatorconfig.json" //Hay que buscar una forma de hacer una especie de NewRequest("Get", file)

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
