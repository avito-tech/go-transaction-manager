package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(mongo.Session).
// TODO add options.
func NewDefaultFactory(client *mongo.Client) transaction.TrFactory {
	return func(ctx context.Context, s transaction.Settings) (context.Context, transaction.Transaction, error) {
		return NewTransaction(ctx, nil, nil, client)
	}
}
