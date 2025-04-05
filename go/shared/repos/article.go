package repos

import (
	"context"
	"fmt"
	"log"

	"github.com/Kaaaaazuya/rss-commentator/go/shared/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gopkg.in/guregu/null.v3"
)

type ArticleRepo struct {
	BaseRepo
	TableName string
}

type ListArticleParameter struct {
	TargetDate null.String
}

type IAritcleRepo interface {
	Create(ctx context.Context, article *models.Article) error
	UpdateSummary(ctx context.Context, hash, text string) error
	List(ctx context.Context, params ListArticleParameter) ([]*models.Article, error)
}

// NewArticleRepo creates a new article repository.
func NewArticleRepo(dbc *dynamodb.Client) *ArticleRepo {
	return &ArticleRepo{
		BaseRepo: BaseRepo{
			DBClient: dbc,
		},
		TableName: "articles",
	}
}

// Create creates a new article.
func (r *ArticleRepo) Create(ctx context.Context, article *models.Article) error {
	params, err := attributevalue.MarshalList([]interface{}{article.UrlHash, article.Url, article.Title, article.Summary, article.CreatedAt})
	if err != nil {
		return err
	}

	log.Println(&dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("INSERT INTO \"%s\" VALUE {'url_hash': ?,'url': ?, 'title': ?, 'summary': ?, 'created_at': ?}",
				r.TableName)),
		Parameters: params})

	_, err = r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("INSERT INTO \"%s\" VALUE {'url_hash': ?,'url': ?, 'title': ?, 'summary': ?, 'created_at': ?}",
				r.TableName)),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't insert an item with PartiQL. Here's why: %v\n", err)
	}

	return err
}

// UpdateSummary updates the summary of an article.
func (r *ArticleRepo) UpdateSummary(ctx context.Context, hash, text string) error {
	stmt := fmt.Sprintf("UPDATE \"%s\" SET summary = ? WHERE url_hash = ?", r.TableName)
	_, err := r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: text},
			&types.AttributeValueMemberS{Value: hash},
		},
	})
	if err != nil {
		log.Printf("Couldn't update the item. Here's why: %v\n", err)
		return err
	}

	return nil
}

// List retrieves a list of articles from the database.
func (r *ArticleRepo) List(ctx context.Context, params ListArticleParameter) ([]*models.Article, error) {
	stmt := fmt.Sprintf("SELECT * FROM \"%s\"", r.TableName)
	result, err := r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
	})
	if err != nil {
		log.Printf("Couldn't scan the table. Here's why: %v\n", err)
		return nil, err
	}

	var articles []*models.Article
	err = attributevalue.UnmarshalListOfMaps(result.Items, &articles)
	if err != nil {
		log.Printf("Couldn't unmarshal the result items. Here's why: %v\n", err)
		return nil, err
	}

	return articles, nil
}
