package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, cmd command) error {
	err := s.Queries.DeleteUsersTable(context.Background())
	if err != nil {
		log.Fatalf("Error deleting users table - %v", err)
	}
	fmt.Println("Users table deleted successfully")
	return nil
}

func handlerAllUsers(s *state, cmd command) error {
	users, err := s.Queries.AllUsers(context.Background())
	if err != nil {
		log.Fatalf("error showcasing user list - %v", err)
	}
	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}
	}
	return nil
}
