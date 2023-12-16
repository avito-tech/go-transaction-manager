package goredis8

import (
	"context"

	"github.com/go-redis/redis/v8"

	trm "github.com/avito-tech/go-transaction-manager/v2"
)

// NewDefaultFactory creates default trm.Transaction(redis.UniversalClient).
func NewDefaultFactory(db redis.UniversalClient) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, db, s)
	}
}
