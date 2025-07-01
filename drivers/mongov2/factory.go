package mongov2

import (
	"context"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2"
)

// NewDefaultFactory creates default trm.Transaction(mongo.Session).
func NewDefaultFactory(client client) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, s.SessionOpts(), s.TransactionOpts(), client)
	}
}
