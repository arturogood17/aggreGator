package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.flags) != 1 { //Tiene que ser así porque solo están pasando 1 flag de main que es el nombre
		return errors.New("usage: login <username>")
	}
	user, err := s.Queries.GetUser(context.Background(), cmd.flags[0])
	if err != nil {
		log.Fatalf("User does not exist in the database")
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return errors.New("error setting user")
	}
	fmt.Println("User set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.flags) != 1 { //Tiene que ser así porque solo están pasando 1 flag de main que es el nombre
		return errors.New("usage: register <username>")
	}
	_, err := s.Queries.GetUser(context.Background(), cmd.flags[0])
	if err == nil {
		log.Fatalf("User already exists")
	}
	user, err := s.Queries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.flags[0],
	})
	if err != nil {
		log.Fatalf("Error creating user - %v", err)
		os.Exit(1)
	}
	s.cfg.SetUser(user.Name)
	fmt.Println("User created and logged.")
	fmt.Println()
	PrettyPrinting(user)
	return nil
}

func PrettyPrinting(user database.User) {
	fmt.Printf("UserID: %v\n", user.ID)
	fmt.Printf("User created at: %v\n", user.CreatedAt)
	fmt.Printf("User updated at: %v\n", user.UpdatedAt)
	fmt.Printf("User's name: %v\n", user.Name)
}
