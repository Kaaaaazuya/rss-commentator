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

	if err = h.handler(); err != nil {
		h.Logger.Error("Error in handler", zap.Error(err))
		return err
	}

	return nil
}

func (h *Handler) handler() error {
	h.Logger.Info("start")
	defer h.Logger.Info("end")

	h.NotifyClient.SendMessage("test")

	return nil
}
