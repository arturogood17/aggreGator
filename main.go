package main

import (
	"fmt"

	"github.com/arturogood17/aggreGator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(cfg.DbURL)
		fmt.Println(cfg.CurrentUserName)
	}
}
