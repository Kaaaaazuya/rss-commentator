package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer/llm"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer/pkg"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/db"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/models"
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
	ArticleRepo    repos.IAritcleRepo
	TagRepo        repos.ITagRepo
	ArticleTagRepo repos.IArticleTagRepo
	Model          llms.Model
}

func NewHandler(dbc *dynamodb.Client, model llms.Model) *Handler {
	return &Handler{
		ArticleRepo:    repos.NewArticleRepo(dbc),
		TagRepo:        repos.NewTagRepo(dbc),
		ArticleTagRepo: repos.NewArticleTagRepo(dbc),
		Model:          model,
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

	if err = h.handler(ctx); err != nil {
		log.Printf("Error handling request: %v", err)
		return err
	}

	return nil
}

func (h *Handler) handler(ctx context.Context) error {
	log.Printf("summarizer started")
	defer log.Printf("summarizer finished")

	today := time.Now().Format("2006-01-02")
	log.Printf("Target date: %s", today)
	targetDate := null.NewString(today, true)

	articles, err := h.ArticleRepo.List(ctx, repos.ListArticleParameter{TargetDate: targetDate})
	if err != nil {
		log.Printf("Error listing articles: %v", err)
		return err
	}

	if len(articles) == 0 {
		log.Printf("No articles found for date: %s", today)
		return nil
	}
	log.Printf("Found %d articles", len(articles))

	tags, err := h.TagRepo.List(ctx)
	if err != nil {
		log.Printf("Error listing tags: %v", err)
		return err
	}
	// propmpt に渡せるようにテキストに変換する
	tagTexts := make([]string, len(tags))
	for i, tag := range tags {
		tagTexts[i] = fmt.Sprintf("・%s\n", tag.TagName)
	}
	// テキストを結合する
	tagsText := fmt.Sprintf("以下のタグを参考にしてください。\n\n%s", tagTexts)

	for _, article := range articles {
		// すでに要約がある場合はスキップする
		if article.Summary != "" {
			log.Printf("Article already has summary: %s", article.Url)
			continue
		}

		prompt := fmt.Sprintf(
			pkg.TEMPLATE,
			tagsText,
			article.Url,
		)
		completion, err := llms.GenerateFromSinglePrompt(ctx, h.Model, prompt)
		if err != nil {
			log.Fatal(err)
			continue
		}

		summary, err := pkg.ParseResponse(completion)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// 取得した要約を記事に保存する
		err = h.ArticleRepo.UpdateSummary(ctx, article.UrlHash, summary.Summary)
		if err != nil {
			log.Printf("Error updating article summary: %v", err)
			continue
		}

		// タグを保存する
		for _, tag := range summary.Tags {
			if tag.IsNew {
				// TODO:  新規タグの場合は、タグを保存する
				continue
			}
			// タグを保存する
			err = h.ArticleTagRepo.Create(ctx, &models.ArticleTag{
				UrlHash: article.UrlHash,
				TagName: tag.Name,
				Score:   tag.Score,
			})
			if err != nil {
				log.Printf("Error creating article tag: %v", err)
				continue
			}
		}

		// スリープ処理を入れる
		time.Sleep(3 * time.Second)
		break
	}

	return nil
}
