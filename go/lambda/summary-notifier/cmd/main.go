package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summary-notifier/line"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summary-notifier/pkg"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/db"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/repos"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v3"
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
		return
	}

	line, err := line.NewClient()
	if err != nil {
		return
	}

	h := NewHandler(logger, dbc, line)

	lambda.Start(h.Handler)
}

type Handler struct {
	Logger         *zap.Logger
	ArticleRepo    repos.IAritcleRepo
	ArticleTagRepo repos.IArticleTagRepo
	NotifyClient   *line.Client
}

func NewHandler(logger *zap.Logger, dbc *dynamodb.Client, line *line.Client) *Handler {
	return &Handler{
		Logger:       logger,
		ArticleRepo:  repos.NewArticleRepo(dbc),
		ArticleTagRepo: repos.NewArticleTagRepo(dbc),
		NotifyClient: line,
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
		h.Logger.Error("Error in handler", zap.Error(err))
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

	for _, article := range articles {
		ats, err := h.ArticleTagRepo.ListByUrlHash(ctx, article.UrlHash)
		if err != nil {
			h.Logger.Error("Error listing article tags", zap.String("url_hash", article.UrlHash), zap.Error(err))
			continue
		}

		message := line.GenerateNotificationMessage(article, ats)
		h.NotifyClient.SendMessage(message)
	}

	return nil
}
