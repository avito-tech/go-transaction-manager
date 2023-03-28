// Package uow is an implementation of trm.Transaction interface by Transaction for trm.UoW.
package uow

import (
	"context"
	"sync/atomic"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/uow"
)

type Transaction struct {
	uow      trm.UoW
	isActive int64
}

// NewTransaction creates trm.Transaction for sqlx.Tx.
func NewTransaction(ctx context.Context, manager trm.Manager) (context.Context, *Transaction, error) {
	tr := &Transaction{uow: uow.NewUoW(manager), isActive: 1}

	go tr.awaitDone(ctx)

	return ctx, tr, nil
}

func (t *Transaction) awaitDone(ctx context.Context) {
	if ctx.Done() == nil {
		return
	}

	<-ctx.Done()

	t.deactivate()
}

func (t *Transaction) Transaction() interface{} {
	return t.uow
}

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	defer t.deactivate()

	return t.uow.Commit(ctx)
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	return nil
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return atomic.LoadInt64(&t.isActive) == 1
}

func (t *Transaction) deactivate() {
	atomic.SwapInt64(&t.isActive, 0)
}
