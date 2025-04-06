package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/arturogood17/aggreGator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name  string
	flags []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	if f, exists := c.cmds[cmd.name]; !exists {
		return errors.New("this command does not exists")
	} else {
		f(s, cmd)
	}
	return nil
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	s := state{
		cfg: &config,
	}
	mapCommands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	//Register segment
	mapCommands.register("login", handlerLogin)

	//Run segment
	var n string
	var f []string
	if len(os.Args) <= 1 {
		log.Fatal("not enough arguments were provided")
		os.Exit(1)
	}
	n = os.Args[1]
	if len(os.Args) <= 2 {
		log.Fatal("username is required")
		os.Exit(1)
	}
	f = os.Args[2:]
	err = mapCommands.run(&s, command{name: n, flags: f})
	if err != nil {
		log.Fatalf("Error running this command: %v. Error value: %v", os.Args[1], err)
	}
	fmt.Println(s.cfg.CurrentUserName)
}
