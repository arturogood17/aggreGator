package main

import (
	"context"
	"fmt"
	"html"

	"log"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func handlerFeedFuncs(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error handling feed - %v", err)
	}
	fmt.Println(html.UnescapeString(feed.Channel.Title))
	fmt.Println(feed.Channel.Link)
	fmt.Println(html.UnescapeString(feed.Channel.Description))
	for _, item := range feed.Channel.Item {
		fmt.Println(html.UnescapeString(item.Title))
		fmt.Println(item.Link)
		fmt.Println(html.UnescapeString(item.Description))
		fmt.Println(item.PubDate)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.flags) != 2 {
		log.Fatalf("usage: %v <name> <url>", cmd.name)
	}
	feed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   cmd.flags[0],
		Url:    cmd.flags[1],
		UserID: user.ID})
	if err != nil {
		return fmt.Errorf("error creating feed - %v", err)
	}
	_, err = s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error trying to follow feed - %v", err)
	}
	fmt.Println(feed.ID.String())
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID.String())
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feedL, err := s.Queries.FeedList(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feed list -%v", err)
	}
	if len(feedL) == 0 { //Tienes que revisar que la lista no esté vacía. Si lo está, no es un error
		fmt.Println("No feeds found")
		return nil
	}
	for _, feed := range feedL {
		username, err := s.Queries.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting the creator's name for the feed -%v", err)
		}
		fmt.Printf("* Name: %v\n", feed.Name)
		fmt.Printf("* URL: %v\n", feed.Url)
		fmt.Printf("* Created by: %v\n", username.Name)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.flags) != 1 {
		log.Fatalf("usage: %v <url>", cmd.name)
	}
	feed, err := s.Queries.FeedByURL(context.Background(), cmd.flags[0])
	if err != nil {
		return fmt.Errorf("error getting feed to follow feed - %v", err)
	}
	followedF, err := s.Queries.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:     uuid.New(),
			UserID: user.ID,
			FeedID: feed.ID,
		})
	if err != nil {
		return fmt.Errorf("error following feed - %v", err)
	}
	fmt.Printf("* Name: %v\n", followedF.FeedName)
	fmt.Printf("* Name: %v\n", followedF.UserName)
	return nil
}

func handlerFollowingFeeds(s *state, cmd command, user database.User) error {
	followedL, err := s.Queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting list of followed feeds - %v", err)
	}
	if len(followedL) == 0 {
		fmt.Printf("User %v is not following any feeds\n", s.cfg.CurrentUserName)
		return nil
	}
	fmt.Printf("* Followed Feeds of: %v\n", s.cfg.CurrentUserName)
	for _, feed := range followedL {
		fmt.Printf("* Feed Name: %v\n", feed)
	}
	return nil
}
