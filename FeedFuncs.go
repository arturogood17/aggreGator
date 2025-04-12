package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error creating request of feeds - %v", err)
	}
	req.Header.Set("User-Agent", "gator")
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error trying to make request - %v", err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading body of response - %v", err)
	}
	defer res.Body.Close()
	var Feed RSSFeed
	if err := xml.Unmarshal(data, &Feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("error unmarshaling data of response - %v", err)
	}
	return &Feed, nil
}
