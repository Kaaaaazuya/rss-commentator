package pkg

import (
	"encoding/json"
	"fmt"
)

// ArticleSummary は正常な記事要約のレスポンスを表します。
type ArticleSummary struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Tags    []Tag  `json:"tags"`
}

// Tag は記事要約に付与されるタグ情報です。
type Tag struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
	IsNew bool    `json:"is_new"`
}

// ErrorResponse は記事の取得に失敗した場合のレスポンスを表します。
type ErrorResponse struct {
	Error string `json:"error"`
}

// ParseResponse はLLMからのJSONレスポンス文字列をパースし、ArticleSummaryまたはErrorResponseに変換します。
func ParseResponse(responseJSON string) (*ArticleSummary, error) {
	// まずエラーレスポンスかチェック
	var errResp ErrorResponse
	if err := json.Unmarshal([]byte(responseJSON), &errResp); err == nil && errResp.Error != "" {
		return nil, fmt.Errorf("LLM error: %s", errResp.Error)
	}

	// 正常なレスポンスとしてパース
	var summary ArticleSummary
	if err := json.Unmarshal([]byte(responseJSON), &summary); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &summary, nil
}
