package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var rss RSSFeed
	if err := xml.Unmarshal(data, &rss); err != nil {
		return &RSSFeed{}, nil
	}
	UnSTitleCh := html.UnescapeString(rss.Channel.Title)
	rss.Channel.Title = UnSTitleCh
	UnSDesCh := html.UnescapeString(rss.Channel.Description)
	rss.Channel.Title = UnSDesCh

	for _, v := range rss.Channel.Item {
		UnSTitle := html.UnescapeString(v.Title)
		v.Title = UnSTitle
		UnSDescr := html.UnescapeString(v.Description)
		v.Title = UnSDescr
	}
	return &rss, nil
}
