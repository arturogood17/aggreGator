package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func addFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("provide name and url")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
		Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: user.ID, FeedID: feed.ID,
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

func followFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("command usage: follow <url>")
	}
	feed, err := s.db.FeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	followed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	for _, v := range followed {
		fmt.Println(v.FeedName)
		fmt.Println(v.UserName)
	}
	return nil
}

func followedList(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Println("no feed follows found for this user")
		return nil
	}
	fmt.Println(user)
	for _, f := range feeds {
		fmt.Printf("* %s", f.Feed)
	}
	return nil
}

func unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("command usage: unfollow <url>")
	}
	ToUnfollow, err := s.db.FeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	if err := s.db.DeleteFeed(context.Background(), database.DeleteFeedParams{UserID: user.ID, FeedID: ToUnfollow.ID}); err != nil {
		return err
	}
	return nil
}

func browsePosts(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		if val, err := strconv.Atoi(cmd.args[0]); err != nil {
			return err
		} else {
			limit = val
		}
	}
	posts, err := s.db.PostsForUser(context.Background(), database.PostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}
	for _, p := range posts {
		fmt.Printf("* %s\n", p.Title)
		fmt.Printf("* %s\n", p.PublishedAt.Time)
		fmt.Printf("* %s\n", p.Description.String)
	}
	return nil
}
