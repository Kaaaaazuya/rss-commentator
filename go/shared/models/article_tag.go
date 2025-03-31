package models

type ArticleTag struct {
	TagName   string  `json:"tag_name" dynamodbav:"tag_name"`
	ArticleID string  `json:"article_id" dynamodbav:"article_id"`
	Score     float64 `json:"score" dynamodbav:"score"`
}
