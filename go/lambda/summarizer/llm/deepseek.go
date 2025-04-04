package llm

import (
	"github.com/Kaaaaazuya/rss-commentator/go/lambda/summarizer/pkg"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// NewDeepSeekClient は DeepSeek API を利用する langchaingo クライアントを生成します。
func NewDeepSeekClient(cfg pkg.LLMConfig) (llms.Model, error) {
	client, err := openai.New(
		openai.WithToken(cfg.APIKey),     // DeepSeek のアクセストークン
		openai.WithBaseURL(cfg.Endpoint), // DeepSeek のエンドポイントURL
		openai.WithModel(cfg.Model),      // 使用するモデル
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
