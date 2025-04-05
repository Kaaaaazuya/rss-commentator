package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/pkg"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/xmltransformer"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/db"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/models"
	"github.com/Kaaaaazuya/rss-commentator/go/shared/repos"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RSSフィードのXML構造体（必要な要素のみ）
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

// "https://feeds.feedburner.com/blogspot/RLXA"
var Urls = []string{"https://zenn.dev/p/acntechjp/feed"}

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

	h := NewHandler(logger, dbc)

	lambda.Start(h.Handler)
}

type Handler struct {
	Logger         *zap.Logger
	ArticleRepo    repos.IAritcleRepo
	XMLTransformer xmltransformer.Transformer
}

func NewHandler(logger *zap.Logger, dbc *dynamodb.Client) *Handler {
	return &Handler{
		Logger:         logger,
		ArticleRepo:    repos.NewArticleRepo(dbc),
		XMLTransformer: xmltransformer.New(),
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

	for _, url := range Urls {
		resp, err := http.Get(url)
		if err != nil {
			h.Logger.Error("Error fetching URL", zap.String("url", url), zap.Error(err))
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			h.Logger.Error("Error reading response body", zap.String("url", url), zap.Error(err))
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		rss, err := h.XMLTransformer.TransformToRSS(body)
		if err != nil {
			h.Logger.Error("Error transforming XML to RSS", zap.String("url", url), zap.Error(err))
			continue
		}
		var articles []models.Article
		for _, item := range rss.Channel.Items {
			// RSS から取得した当日の記事のみを取得
			pubDateParsed, err := time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				// パースに失敗した場合の処理
				h.Logger.Info("Error parsing PubDate", zap.String("pubDate", item.PubDate), zap.Error(err))
				continue
			}

			// JST のタイムゾーンを設定
			jst := time.FixedZone("Asia/Tokyo", 9*60*60)
			pubDateJST := pubDateParsed.In(jst)
			nowJST := time.Now().In(jst)

			// 年、月、日だけを比較
			if pubDateJST.Year() != nowJST.Year() ||
				pubDateJST.Month() != nowJST.Month() ||
				pubDateJST.Day() != nowJST.Day() {
				continue
			}

			hash := sha256.Sum256([]byte(item.Link))
			hashString := hex.EncodeToString(hash[:])

			// 取得した記事情報を models.Article に変換
			// Summary は要約後に設定するため空のままとする
			article := models.Article{
				UrlHash:   hashString,
				Url:       item.Link,
				Title:     item.Title,
				Summary:   "",
				CreatedAt: pubDateJST.Format(time.RFC3339),
			}
			articles = append(articles, article)
		}

		// 取得した記事情報を表示
		for _, art := range articles {
			// DB に保存
			err := h.ArticleRepo.Create(context.Background(), &art)
			if err != nil {
				h.Logger.Error("Error saving article to DB", zap.String("url", art.Url), zap.Error(err))
				continue
			}
		}
	}

	return nil
}
