package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/db"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/models"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/pkg"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/repos"
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/xmltransformer"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
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
	XMLTransformer xmltransformer.Transformer
}

func NewHandler(dbc *dynamodb.Client) *Handler {
	return &Handler{
		DBClient:       dbc,
		ArticleRepo:    repos.NewArticleRepo(dbc),
		XMLTransformer: xmltransformer.New(),
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
	log.Printf("rss-fetcher started")
	defer log.Printf("rss-fetcher finished")

	for _, url := range Urls {
		log.Printf("Fetching RSS feed from %s", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching %s: %v", url, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body from %s: %v", url, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		rss, err := h.XMLTransformer.TransformToRSS(body)
		if err != nil {
			log.Printf("Error transforming XML to RSS: %v", err)
		}
		var articles []models.Article
		for _, item := range rss.Channel.Items {
			// RSS から取得した当日の記事のみを取得
			pubDateParsed, err := time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				// パースに失敗した場合の処理
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

			// 取得した記事情報を models.Article に変換
			// Summary は要約後に設定するため空のままとする
			article := models.Article{
				Id:        uuid.New(),
				Url:       item.Link,
				Title:     item.Title,
				Summary:   "",
				CreatedAt: pubDateJST.Format(time.RFC3339),
			}
			articles = append(articles, article)
		}

		// 例として、取得した記事情報を表示
		for _, art := range articles {
			fmt.Printf("Article: %+v\n", art)
		}
	}

	return nil
}
