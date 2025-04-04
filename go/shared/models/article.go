package models

// Article represents an article
type Article struct {
	UrlHash   string `json:"urlHash" dynamodbav:"urlHash"`
	Url       string `json:"url" dynamodbav:"url"`
	Title     string `json:"title" dynamodbav:"title"`
	Summary   string `json:"summary" dynamodbav:"summary"`
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`
}
