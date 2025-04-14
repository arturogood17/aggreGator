package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state, ctx context.Context) {
	feeds, err := s.Queries.FeedList(ctx)
	if err != nil {
		fmt.Println("error not getting feeds")
		return
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds in to fetch from")
		return
	}
	feed, err := s.Queries.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("error getting feed - %v\n", err)
		return
	}
	fmt.Println("Found feed to fetch!")

	if err = s.Queries.MarkedAsFetched(ctx, feed.ID); err != nil {
		fmt.Printf("error marking feed as fetch - %v\n", err)
		return
	}
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		fmt.Printf("error fetching feed - %v\n", err)
		return
	}
	for _, post := range rssFeed.Channel.Item {
		published, err := time.Parse(time.RFC1123, post.PubDate)
		if err != nil {
			fmt.Printf("Error trying to created published date - %v\n", err)
			return
		}
		if _, err := s.Queries.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			Title:       post.Title,
			Url:         post.Link,
			Description: sql.NullString{String: post.Description, Valid: true},
			PublishedAt: published,
			FeedID:      feed.ID,
		}); err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			} else {
				fmt.Printf("Error saving post in database - %v\n", err)
				return
			}
		}
	}
	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
