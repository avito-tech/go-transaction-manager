package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewFactory creates default transaction.Transaction(sqlx.Tx).
func NewFactory(db *sql.DB) transaction.TrFactory {
	return func(ctx context.Context) (transaction.Transaction, error) {
		return NewTransaction(ctx, transaction.NewSavePoint(), nil, db)
	}
}
