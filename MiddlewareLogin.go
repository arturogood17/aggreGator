package main

import (
	"context"
	"fmt"

	"github.com/arturogood17/aggreGator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		validUser, err := s.Queries.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting user - %v", err)
		}
		return handler(s, cmd, validUser) //tienes que llamar la función que te pasaron de parámetro después de hacer la validación
		// puedes devolver la función porque la función de adentro devuelve un error
	}
}
