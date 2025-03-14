package config

import (
	"encoding/json"
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
	file += "~/.gatorconfig.json" //Hay que buscar una forma de hacer una especie de NewRequest("Get", file)
	var cfg Config
	json.Unmarshal(file, &cfg)
}
