package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Client struct {
	accessToken string
}

func NewClient() (*Client, error) {
	accessToken := os.Getenv("LINE_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("LINE_ACCESS_TOKEN is not set")
	}
	return &Client{accessToken: accessToken}, nil
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (c *Client) SendMessage(message string) error {
	url := "https://api.line.me/v2/bot/message/broadcast"

	msg := Message{
		Type: "text",
		Text: message,
	}

	payload := map[string][]Message{
		"messages": {msg},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
