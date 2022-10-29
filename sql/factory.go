package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(sql.Tx).
func NewDefaultFactory(db *sql.DB) transaction.TrFactory {
	return NewFactory(db, NewSavePoint())
}

// NewFactory creates transaction.Transaction(sql.Tx).
func NewFactory(db *sql.DB, sp SavePoint) transaction.TrFactory {
	return func(ctx context.Context, trms transaction.Settings) (context.Context, transaction.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, sp, s.TxOpts(), db)
	}
}
