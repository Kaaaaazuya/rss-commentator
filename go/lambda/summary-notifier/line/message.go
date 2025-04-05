package line

import (
	"fmt"
	"strings"

	"github.com/Kaaaaazuya/rss-commentator/go/shared/models"
)

// GenerateNotificationMessage は記事とタグ情報からLINE通知用のメッセージを生成します。
func GenerateNotificationMessage(a *models.Article, ats []*models.ArticleTag) string {
	// タグ情報を "#タグ名" の形式に整形
	var tagList []string
	for _, at := range ats {
		tagList = append(tagList, fmt.Sprintf("#%s", at.TagName))
	}
	tagsStr := strings.Join(tagList, " ")

	// 通知メッセージの組み立て
	message := fmt.Sprintf(
		"【新着記事】%s\n\n%s\n\n詳細はこちら: %s\n%s",
		a.Title,
		a.Summary,
		a.Url,
		tagsStr,
	)

	return message
}
