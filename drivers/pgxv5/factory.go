//go:build go1.19
// +build go1.19

package pgxv5

import (
	"context"

	trm "github.com/avito-tech/go-transaction-manager/v2"
)

// NewDefaultFactory creates default trm.Transaction(pgx.Tx).
func NewDefaultFactory(db Transactional) trm.TrFactory {
	return NewFactory(db)
}

// NewFactory creates trm.Transaction(pgx.Tx).
func NewFactory(db Transactional) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, s.TxOpts(), db)
	}
}
