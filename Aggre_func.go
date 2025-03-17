package main

import (
	"context"
	"fmt"
)

func Aggregation(s *state, cmd command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rss.Channel.Title)
	fmt.Println(rss.Channel.Link)
	fmt.Println(rss.Channel.Description)
	for _, v := range rss.Channel.Item {
		fmt.Println(v.Title)
		fmt.Println(v.Link)
		fmt.Println(v.Description)
		fmt.Println(v.PubDate)
	}
	return nil
}
