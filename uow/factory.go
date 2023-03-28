package uow

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
)

func NewDefaultFactory(manager trm.Manager) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		return NewTransaction(ctx, manager)
	}
}
