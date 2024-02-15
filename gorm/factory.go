//go:build go1.16
// +build go1.16

package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/avito-tech/go-transaction-manager/trm"
)

// NewDefaultFactory creates default trm.Transaction(gorm.DB).
// Factory rewrites DisableNestedTransaction in gorm.Config with Propagation in trm.Settings.
func NewDefaultFactory(db *gorm.DB) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		db.Config.DisableNestedTransaction = trms.Propagation() != trm.PropagationNested

		return NewTransaction(ctx, s.TxOpts(), db)
	}
}
