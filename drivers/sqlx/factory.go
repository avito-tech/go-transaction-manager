package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"

	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
)

// NewDefaultFactory creates default trm.Transaction(sqlx.Tx).
func NewDefaultFactory(db *sqlx.DB) trm.TrFactory {
	return NewFactory(db, trmsql.NewSavePoint())
}

// NewFactory creates trm.Transaction(sql.Tx).
func NewFactory(db *sqlx.DB, sp trmsql.SavePoint) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(trmsql.Settings)

		return NewTransaction(ctx, sp, s.TxOpts(), db)
	}
}
