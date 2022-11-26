package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/avito-tech/go-transaction-manager/trm"
)

// NewDefaultFactory creates default trm.Transaction(redis.UniversalClient).
func NewDefaultFactory(db redis.UniversalClient) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, db, s)
	}
}
