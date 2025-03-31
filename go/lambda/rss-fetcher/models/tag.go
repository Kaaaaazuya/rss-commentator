package models

type Tag struct {
	Name      string `json:"name" dynamodbav:"name"`
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`
}
