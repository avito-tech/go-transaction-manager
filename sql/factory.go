package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/trm"
)

// NewDefaultFactory creates default trm.Transaction(sql.Tx).
func NewDefaultFactory(db *sql.DB) trm.TrFactory {
	return NewFactory(db, NewSavePoint())
}

// NewFactory creates trm.Transaction(sql.Tx).
func NewFactory(db *sql.DB, sp SavePoint) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, sp, s.TxOpts(), db)
	}
}
