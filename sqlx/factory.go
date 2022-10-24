package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(sqlx.Tx).
func NewDefaultFactory(db *sqlx.DB) transaction.TrFactory {
	return func(ctx context.Context) (transaction.Transaction, error) {
		return NewTransaction(ctx, transaction.NewSavePoint(), nil, db)
	}
}
