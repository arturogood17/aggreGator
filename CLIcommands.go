package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registeredcmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredcmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	n := cmd.name
	f, ok := c.registeredcmds[n]
	if !ok {
		return fmt.Errorf("command not found")
	}
	return f(s, cmd)
}
