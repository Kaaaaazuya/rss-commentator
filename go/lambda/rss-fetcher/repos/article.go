package repos

import (
	"context"
	"fmt"
	"log"

	"github.com/Kaaaaazuya/rss-commentator/go/lambda/rss-fetcher/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ArticleRepo struct {
	BaseRepo
	TableName string
}

type IAritcleRepo interface {
	Create(ctx context.Context, article *models.Article) error
}

// NewArticleRepo creates a new article repository.
func NewArticleRepo(dbc *dynamodb.Client) *ArticleRepo {
	return &ArticleRepo{
		BaseRepo: BaseRepo{
			DBClient: dbc,
		},
		TableName: "Articles",
	}
}

// Create creates a new article.
func (r *ArticleRepo) Create(ctx context.Context, article *models.Article) error {
	params, err := attributevalue.MarshalList([]interface{}{article.Url, article.Title, article.Summary, article.CreatedAt})
	if err != nil {
		return err
	}

	_, err = r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("INSERT INTO \"%s\" VALUE {'id': ?, 'url': ?, 'title': ?, 'summary': ?, 'tags': ?, 'createdAt': ?}",
				r.TableName)),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't insert an item with PartiQL. Here's why: %v\n", err)
	}

	return err
}
