package models

import (
	"github.com/google/uuid"
)

// Article represents an article
type Article struct {
	Id        uuid.UUID `json:"id" dynamodbav:"id"`
	Url       string    `json:"url" dynamodbav:"url"`
	Title     string    `json:"title" dynamodbav:"title"`
	Summary   string    `json:"summary" dynamodbav:"summary"`
	CreatedAt string    `json:"createdAt" dynamodbav:"createdAt"`
}
