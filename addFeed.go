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
	fmt.Println()
	feedPrinting(feed)
	fmt.Println()
	return nil
}

func feedPrinting(f database.Feed) {
	fmt.Println(f.ID)
	fmt.Println(f.CreatedAt)
	fmt.Println(f.UpdatedAt)
	fmt.Println(f.Name)
	fmt.Println(f.Url)
	fmt.Println(f.UserID)
}

func feedList(s *state, cmd command) error {
	feed_list, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feed_list {
		fmt.Println(f.Name)
		fmt.Println(f.Url)
		user, err := s.db.GetUserByID(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Println(user.Name)
	}
	return nil
}
