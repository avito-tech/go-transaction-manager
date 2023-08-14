package pgx

import (
	"context"

	trmsql "github.com/avito-tech/go-transaction-manager/sql"

	"github.com/avito-tech/go-transaction-manager/trm"
)

// NewDefaultFactory creates default trm.Transaction(pgx.Tx).
func NewDefaultFactory(db Transactional) trm.TrFactory {
	return NewFactory(db, trmsql.NewSavePoint())
}

// NewFactory creates trm.Transaction(pgx.Tx).
func NewFactory(db Transactional, sp trmsql.SavePoint) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, sp, s.TxOpts(), db)
	}
}
