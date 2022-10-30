package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"

	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(sqlx.Tx).
func NewDefaultFactory(db *sqlx.DB) transaction.TrFactory {
	return NewFactory(db, trmsql.NewSavePoint())
}

// NewFactory creates transaction.Transaction(sql.Tx).
func NewFactory(db *sqlx.DB, sp trmsql.SavePoint) transaction.TrFactory {
	return func(ctx context.Context, trms transaction.Settings) (context.Context, transaction.Transaction, error) {
		s, _ := trms.(trmsql.Settings)

		return NewTransaction(ctx, sp, s.TxOpts(), db)
	}
}
