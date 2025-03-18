package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func addFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("provide name and url")
	}
	currentUser := s.cfg.CurrentUserName
	UI, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return err
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
		Name: cmd.args[0], Url: cmd.args[1], UserID: UI.ID})
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
