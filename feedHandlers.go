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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.flags) != 2 {
		log.Fatalf("usage: %v <name> <url>", cmd.name)
	}
	user, err := s.Queries.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user to create feed - %v", err)
	}
	feed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   cmd.flags[0],
		Url:    cmd.flags[1],
		UserID: user.ID})
	if err != nil {
		return fmt.Errorf("error creating feed - %v", err)
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
