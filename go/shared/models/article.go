package models

// Article represents an article
type Article struct {
	UrlHash   string `json:"url_hash" dynamodbav:"url_hash"`
	Url       string `json:"url" dynamodbav:"url"`
	Title     string `json:"title" dynamodbav:"title"`
	Summary   string `json:"summary" dynamodbav:"summary"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
}
