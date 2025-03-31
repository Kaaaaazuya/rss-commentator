package repos

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type BaseRepo struct {
	DBClient *dynamodb.Client
}

type ctxKey string

const txKey ctxKey = "tx"

// Client returns the dynamodb client.
// Transactional client is returned if it is present in the context.
func (r *BaseRepo) Client(ctx context.Context) *dynamodb.Client {
	if tx, ok := ctx.Value(txKey).(*dynamodb.Client); ok && tx != nil {
		return tx
	}
	return r.DBClient
}
