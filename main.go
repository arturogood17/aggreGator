package main

import (
	"fmt"
	"log"

	"github.com/arturogood17/aggreGator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = cfg.SetUser("nestor")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(cfg)
}
