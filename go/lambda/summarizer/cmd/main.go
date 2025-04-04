package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer/llm"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer/pkg"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/db"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/repos"
	"github.com/tmc/langchaingo/llms"
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
		log.Printf("Error creating DynamoDB client: %v", err)
		return
	}

	model, err := llm.NewDeepSeekClient(*c.LLM)
	if err != nil {
		log.Printf("Error creating LLM client: %v", err)
		return
	}

	h := NewHandler(dbc, model)

	lambda.Start(h.Handler)
}

type Handler struct {
	ArticleRepo repos.IAritcleRepo
	Model       llms.Model
}

func NewHandler(dbc *dynamodb.Client, model llms.Model) *Handler {
	return &Handler{
		ArticleRepo: repos.NewArticleRepo(dbc),
		Model:       model,
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
		if err != nil {
			log.Printf("Error creating summarizer: %v", err)
			continue
		}
		prompt := fmt.Sprintf(
			"以下の記事を読み込み、主要なポイント、結論、背景情報を踏まえた上で、3〜5文の簡潔で分かりやすい要約文を生成してください。\n\n記事のURL: %s",
			article.Url,
		)
		completion, err := llms.GenerateFromSinglePrompt(context.Background(), h.Model, prompt)
		if err != nil {
			log.Fatal(err)
		}

		// 結果を表示する
		fmt.Println("===answet===")
		fmt.Println(completion)
		fmt.Println("============")
		break
	}

	return nil
}
