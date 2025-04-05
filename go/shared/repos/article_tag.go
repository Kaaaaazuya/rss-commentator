package repos

import (
	"context"
	"fmt"

	"github.com/Kaaaaazuya/rss-commentator/go/shared/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ArticleTagRepo struct {
	BaseRepo
	TableName string
}

type IArticleTagRepo interface {
	Create(ctx context.Context, articleTag *models.ArticleTag) error
	ListByUrlHash(ctx context.Context, urlHash string) ([]*models.ArticleTag, error)
}

func NewArticleTagRepo(dbc *dynamodb.Client) *ArticleTagRepo {
	return &ArticleTagRepo{
		BaseRepo: BaseRepo{
			DBClient: dbc,
		},
		TableName: "article_tags",
	}
}

// Create creates a new article tag.
func (r *ArticleTagRepo) Create(ctx context.Context, articleTag *models.ArticleTag) error {
	params, err := attributevalue.MarshalList([]interface{}{articleTag.UrlHash, articleTag.TagName, articleTag.Score})
	if err != nil {
		return err
	}

	_, err = r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("INSERT INTO \"%s\" VALUE {'url_hash': ?, 'tag_name': ?, 'score': ?}",
				r.TableName)),
		Parameters: params,
	})
	if err != nil {
		return err
	}

	return nil
}

// ListByUrlHash retrieves article tags by URL hash.
func (r *ArticleTagRepo) ListByUrlHash(ctx context.Context, urlHash string) ([]*models.ArticleTag, error) {
	params, err := attributevalue.MarshalList([]interface{}{urlHash})
	if err != nil {
		return nil, err
	}

	res, err := r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%s\" WHERE url_hash = ?", r.TableName)),
		Parameters: params,
	})
	if err != nil {
		return nil, err
	}

	var articleTags []*models.ArticleTag
	err = attributevalue.UnmarshalListOfMaps(res.Items, &articleTags)
	if err != nil {
		return nil, err
	}

	return articleTags, nil
}
