package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"log"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func handlerFeedFuncs(s *state, cmd command) error {
	if len(cmd.flags) != 1 {
		log.Fatalf("usage: %v <name> <time_between_reqs>", cmd.name)
	}
	time_between_reqs, err := time.ParseDuration(cmd.flags[0])
	if err != nil {
		return fmt.Errorf("error parsing time into time duration - %v", err)
	}

	fmt.Printf("Collecting feed every %v\n", time_between_reqs)

	tick := time.NewTicker(time_between_reqs)

	for ; ; <-tick.C {
		fmt.Println("Collecting...")
		scrapeFeeds(s, context.Background())
	} //no necesitas devolver nada aunque tenga que devolver un error
} //go entiende que es un bucle sin fin

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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.flags) != 1 {
		log.Fatalf("usage: %v <url>", cmd.name)
	}
	feed, err := s.Queries.FeedByURL(context.Background(), cmd.flags[0])
	if err != nil {
		return fmt.Errorf("error getting URL feed to try to delete - %v", err)
	}
	if err := s.Queries.DeleteFeedFollow(context.Background(),
		database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID}); err != nil {
		return fmt.Errorf("error deleting feed - %v", err)
	}
	fmt.Printf("Feed %v successfully deleted\n", feed.Name)
	return nil
}

func handlerBrowsingPosts(s *state, cmd command, user database.User) error {
	if len(cmd.flags) > 1 {
		return fmt.Errorf("usage: %v <limit_of_posts>", cmd.name)
	}
	limit := 2
	if len(cmd.flags) != 0 {
		var err error
		limit, err = strconv.Atoi(cmd.flags[0])
		if err != nil {
			return fmt.Errorf("error converting string into int - %v", err)
		}
	}
	posts, err := s.Queries.GetPostsForUsers(context.Background(), database.GetPostsForUsersParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error fetching posts to show - %v", err)
	}
	for _, p := range posts {
		fmt.Printf("* ID: %v\n", p.ID)
		fmt.Printf("* Title: %v\n", p.Title)
		fmt.Printf("* Description: %v\n", p.Description.String)
		fmt.Printf("* Published at: %v\n", p.PublishedAt)
		fmt.Println()
	}
	return nil
}
