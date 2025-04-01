package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer/pkg"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/db"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/repos"
	"gopkg.in/guregu/null.v3"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	ctx := context.Background()
	c, err := pkg.LoadConfig(ctx)
	if err != nil {
		return
	}

	dbc, err := db.NewClient(ctx, c.AWS)
	if err != nil {
		return
	}

	h := NewHandler(dbc)

	lambda.Start(h.Handler)
}

type Handler struct {
	DBClient       *dynamodb.Client
	ArticleRepo    repos.IAritcleRepo
}

func NewHandler(dbc *dynamodb.Client) *Handler {
	return &Handler{
		ArticleRepo:    repos.NewArticleRepo(dbc),
	}
}

func (h *Handler) Handler(ctx context.Context) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	if err = h.handler(); err != nil {
		log.Printf("Error handling request: %v", err)
		return err
	}

	return nil
}

func (h *Handler) handler() error {
	log.Printf("summarizer started")
	defer log.Printf("summarizer finished")

	today := time.Now().Format("2006-01-02")
	log.Printf("Target date: %s", today)
	targetDate := null.NewString(today, true)

	articles, err := h.ArticleRepo.List(context.Background(), repos.ListArticleParameter{TargetDate: targetDate})
	if err != nil {
		log.Printf("Error listing articles: %v", err)
	}

	for _, article := range articles {
		log.Printf("Article: %v", article)
	}

	return nil
}