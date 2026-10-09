package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mariolazzari/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("get next feed: %w", err)
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("fetch feed %q: %w", feed.Url, err)
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			log.Printf("parse published date for %q: %v", item.Title, err)
			continue
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("save post %q: %v", item.Title, err)
		}
	}

	_, err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("mark feed as fetched: %w", err)
	}

	return nil
}
