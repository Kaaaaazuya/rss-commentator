package repos

import (
	"context"
	"fmt"
	"log"

	"github.com/Kaaaaazuya/rss-commentator/go/shared/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TagRepo struct {
	BaseRepo
	TableName string
}

type ITagRepo interface {
	List(ctx context.Context) ([]*models.Tag, error)
}

func NewTagRepo(dbc *dynamodb.Client) *TagRepo {
	return &TagRepo{
		BaseRepo: BaseRepo{
			DBClient: dbc,
		},
		TableName: "tags",
	}
}

// List retrieves a list of tags from the database.
func (r *TagRepo) List(ctx context.Context) ([]*models.Tag, error) {
	stmt := fmt.Sprintf("SELECT * FROM \"%s\"", r.TableName)
	result, err := r.Client(ctx).ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: &stmt,
	})
	if err != nil {
		log.Printf("Couldn't scan the table. Here's why: %v\n", err)
		return nil, err
	}

	var tags []*models.Tag
	err = attributevalue.UnmarshalListOfMaps(result.Items, &tags)
	if err != nil {
		log.Printf("Couldn't unmarshal the result items. Here's why: %v\n", err)
		return nil, err
	}

	return tags, nil
}
