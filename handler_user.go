package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: cmd <name>")
	}
	username := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		fmt.Println("User login not in database")
		os.Exit(1)
		return err
	}
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User set successfully.")
	return nil
}

func registerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: cmd <name>")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		fmt.Printf("Error: User '%s' already exists\n", cmd.args[0])
		os.Exit(1)
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0]})
	if err != nil {
		return err
	}
	s.cfg.SetUser(user.Name)
	fmt.Println(user.ID)
	fmt.Println(user.CreatedAt)
	fmt.Println(user.UpdatedAt)
	fmt.Println(user.Name)
	return nil
}

func dbDelete(s *state, cmd command) error {
	err := s.db.Delete(context.Background())
	if err != nil {
		fmt.Println("Failed to delete the db.")
		return err
	}
	fmt.Println("Db deleted.")
	return nil
}

func getUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
			continue
		}
		fmt.Printf("* %v\n", user)
	}
	return nil
}
