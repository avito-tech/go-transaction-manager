package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(sql.Tx).
// TODO add options.
func NewDefaultFactory(db *sql.DB) transaction.TrFactory {
	return func(ctx context.Context, s transaction.Settings) (context.Context, transaction.Transaction, error) {
		return NewTransaction(ctx, transaction.NewSavePoint(), nil, db)
	}
}
