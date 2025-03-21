package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/arturogood17/aggreGator/internal/database"
	"github.com/google/uuid"
)

func Aggregation(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}
	aggTime, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	log.Printf("Collecting feeds every: %v", aggTime)

	ticker := time.NewTicker(aggTime)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feedToScrape, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	scrapeFeed(s.db, feedToScrape)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	if _, err := db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(feed.Url)
	fetched, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, v := range fetched.Channel.Item {
		published, err := time.Parse(time.RFC1123, v.PubDate)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       v.Title,
			Url:         v.Link,
			Description: sql.NullString{String: v.Description},
			PublishedAt: sql.NullTime{Time: published},
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post")
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(fetched.Channel.Item))
}
