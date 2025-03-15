package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: cmd <name>")
	}
	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != err {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User set successfully.")
	return nil
}
