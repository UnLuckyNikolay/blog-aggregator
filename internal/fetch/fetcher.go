package fetch

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func fetchFeeds(s *state.State) {
	feeds, err := s.Db.GetFeedsToFetch(context.Background())
	if err != nil {
		fmt.Printf("Error fetching feed list: %s\n", err)
	}

	for _, nextFeed := range feeds {
		err = s.Db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			Time:   time.Now(),
			FeedID: nextFeed.ID,
		})
		if err != nil {
			fmt.Printf("Error making feed fetched: %s\n", err)
		}

		req, err := http.NewRequest("GET", nextFeed.Url, nil)
		if err != nil {
			fmt.Printf("Error making a request: %s\n", err)
		}

		req.Header.Set("User-Agent", "Gator")

		fmt.Printf("Sending a 'GET' request to '%s' (%s)...\n", nextFeed.Name, nextFeed.Url)
		res, err := s.HttpClient.Do(req)
		if err != nil {
			fmt.Printf("Error sending a request: %s\n", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading a request: %s\n", err)
		}

		var feed RSSFeed
		err = xml.Unmarshal(body, &feed)
		if err != nil {
			fmt.Printf("Error unmarshaling a request: %s\n", err)
		}

		feed = unescapeFeed(feed)

		fmt.Printf("New posts:\n")
		for _, item := range feed.Channel.Item {
			post, err := s.Db.AddPost(context.Background(), database.AddPostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				Title:       toNullString(item.Title),
				Url:         item.Link,
				Description: toNullString(item.Description),
				PublishedAt: toNullTime(item.PubDate),
				FeedID:      nextFeed.ID,
			},
			)
			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) {
					if pqErr.Code == "23505" { // Unique URL violation
						continue
					}
					fmt.Printf("fetchFeeds: db error (%s): %s\n", pqErr.Code, err)
				}
				fmt.Printf("fetchFeeds: error: %s\n", err)
				continue
			}
			fmt.Printf("> `%s`\n", post.Title.String)
		}
		fmt.Println()
	}
}

func unescapeFeed(feed RSSFeed) RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item { // Do not use `for _, item` here
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return feed
}
