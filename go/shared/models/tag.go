package models

type Tag struct {
	TagName   string `json:"tag_name" dynamodbav:"tag_name"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
}
