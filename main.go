package main

import (
	"log"
	"os"

	"github.com/arturogood17/aggreGator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	programS := &state{
		cfg: &cfg,
	}
	comnds := commands{
		registeredcmds: make(map[string]func(*state, command) error),
	}
	comnds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err := comnds.run(programS, command{name: cmdName, args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
