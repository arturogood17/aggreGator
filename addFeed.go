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
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: UI.ID, FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("Feed created successfully:")
	feedPrinting(feed)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	fmt.Println("=====================================")
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

func followFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("command usage: follow <url>")
	}
	url := cmd.args[0]
	un := s.cfg.CurrentUserName
	current_user, err := s.db.GetUser(context.Background(), un)
	if err != nil {
		return err
	}
	feed, err := s.db.FeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	followed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: current_user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	for _, v := range followed {
		fmt.Println(v.FeedName)
		fmt.Println(v.UserName)
	}
	return nil
}

func followedList(s *state, cmd command) error {
	user := s.cfg.CurrentUserName
	us, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		return err
	}
	user_id := us.ID
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user_id)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		return errors.New("No feed follows found for this user.")
	}
	fmt.Println(user)
	for _, f := range feeds {
		fmt.Printf("* %s", f.Feed)
	}
	return nil
}
