package mongo

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(mongo.Session).
func NewDefaultFactory(client client) transaction.TrFactory {
	return func(ctx context.Context, trms transaction.Settings) (context.Context, transaction.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, s.SessionOpts(), s.TransactionOpts(), client)
	}
}
