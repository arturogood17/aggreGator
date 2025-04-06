package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.flags) == 0 {
		return errors.New("usage: login <username>")
	}
	err := s.cfg.SetUser(cmd.flags[0])
	if err != nil {
		return errors.New("error setting user")
	}
	fmt.Println("User set")
	return nil
}
