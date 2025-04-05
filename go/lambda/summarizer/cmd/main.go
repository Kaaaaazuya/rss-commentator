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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v3"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cnf := zap.NewProductionConfig()
	cnf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger := zap.Must(cnf.Build())
	defer logger.Sync()

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

	h := NewHandler(logger, dbc, model)

	lambda.Start(h.Handler)
}

type Handler struct {
	Logger         *zap.Logger
	ArticleRepo    repos.IAritcleRepo
	TagRepo        repos.ITagRepo
	ArticleTagRepo repos.IArticleTagRepo
	Model          llms.Model
}

func NewHandler(logger *zap.Logger, dbc *dynamodb.Client, model llms.Model) *Handler {
	return &Handler{
		Logger:         logger,
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
			h.Logger.Error("Recovered from panic", zap.Any("error", r))
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	if err = h.handler(ctx); err != nil {
		h.Logger.Error("Error handling request", zap.Error(err))
		return err
	}

	return nil
}

func (h *Handler) handler(ctx context.Context) error {
	h.Logger.Info("start")
	defer h.Logger.Info("end")

	today := time.Now().Format("2006-01-02")
	targetDate := null.NewString(today, true)

	articles, err := h.ArticleRepo.List(ctx, repos.ListArticleParameter{TargetDate: targetDate})
	if err != nil {
		h.Logger.Error("Error listing articles", zap.Error(err))
		return err
	}

	if len(articles) == 0 {
		h.Logger.Info("No articles found for date", zap.String("date", today))
		return nil
	}
	log.Printf("Found %d articles", len(articles))

	tags, err := h.TagRepo.List(ctx)
	if err != nil {
		h.Logger.Error("Error listing tags", zap.Error(err))
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
			h.Logger.Info("Article already has summary", zap.String("url", article.Url))
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
			h.Logger.Error("Error parsing response", zap.Error(err))
			continue
		}

		// 取得した要約を記事に保存する
		err = h.ArticleRepo.UpdateSummary(ctx, article.UrlHash, summary.Summary)
		if err != nil {
			h.Logger.Error("Error updating article summary", zap.Error(err))
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
				h.Logger.Error("Error creating article tag", zap.Error(err))
				continue
			}
		}

		// スリープ処理を入れる
		time.Sleep(3 * time.Second)
		break
	}

	return nil
}
