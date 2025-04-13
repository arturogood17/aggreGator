package main

import (
	"context"
	"fmt"
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
		fmt.Printf("* Title: %v\n", post.Title)
	}
	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
