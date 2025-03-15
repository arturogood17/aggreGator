package main

import (
	"fmt"

	"github.com/arturogood17/aggreGator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	args []string
	name string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	n := cmd.name
	f, ok := c.cmds[n]
	if !ok {
		return fmt.Errorf("command not found")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username requiered")
	}
	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != err {
		return fmt.Errorf("error setting user")
	}
	fmt.Println("User set.")
	return nil
}
