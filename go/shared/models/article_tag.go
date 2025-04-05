package models

type ArticleTag struct {
	TagName string  `json:"tag_name" dynamodbav:"tag_name"`
	UrlHash string  `json:"url_hash" dynamodbav:"url_hash"`
	Score   float64 `json:"score" dynamodbav:"score"`
}
