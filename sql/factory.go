package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(sqlx.Tx).
func NewDefaultFactory(db *sql.DB) transaction.TrFactory {
	return func(ctx context.Context) (transaction.Transaction, error) {
		return NewTransaction(ctx, transaction.NewSavePoint(), nil, db)
	}
}
