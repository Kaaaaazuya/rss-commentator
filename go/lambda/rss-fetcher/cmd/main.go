package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Feed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
}

var Urls = []string{"https://feeds.feedburner.com/blogspot/RLXA"}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	log.Printf("rss-fetcher started")

	for _, url := range Urls{
		log.Printf("Fetching RSS feed from %s", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching %s: %v", url, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body from %s: %v", url, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		var rss Feed
		err = xml.Unmarshal(body, &rss)
		if err != nil {
			log.Printf("Error parsing RSS feed from %s: %v", url, err)
		} else {
			log.Printf("Fetched feed: %s", rss.Channel.Title)
			
			// Generate record for DynamoDB logging
			record := struct {
				FeedURL   string   `json:"feedURL"`
				FeedTitle string   `json:"feedTitle"`
				FetchedAt string   `json:"fetchedAt"`
				Items     []string `json:"items"`
			}{
				FeedURL:   url,
				FeedTitle: rss.Channel.Title,
				FetchedAt: time.Now().Format(time.RFC3339),
				Items:     []string{},
			}

			recordJSON, err := json.Marshal(record)
			if err != nil {
				log.Printf("Error marshaling record for %s: %v", url, err)
			} else {
				log.Printf("DB Record: %s", recordJSON)
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
