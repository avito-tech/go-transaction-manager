//go:build go1.16
// +build go1.16

package pgxv5

import (
	"context"
	"sync/atomic"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/jackc/pgx/v5"
)

// Transaction is trm.Transaction for pgx.Tx.
type Transaction struct {
	tx       pgx.Tx
	isActive int64
}

// NewTransaction creates trm.Transaction for pgx.Tx.
func NewTransaction(
	ctx context.Context,
	txOptions *pgx.TxOptions,
	db Transactional,
) (context.Context, *Transaction, error) {
	var opts pgx.TxOptions
	if txOptions != nil {
		opts = *txOptions
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{
		tx:       tx,
		isActive: 1,
	}

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

// Transaction returns the real transaction pgx.Tx.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction by save point.
func (t *Transaction) Begin(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) { //nolint:ireturn,nolintlint
	tx, err := t.tx.Begin(ctx)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{
		tx:       tx,
		isActive: 1,
	}

	return ctx, tr, nil
}

// Commit the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	defer t.deactivate()

	return t.tx.Commit(ctx)
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	defer t.deactivate()

	return t.tx.Rollback(ctx)
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return atomic.LoadInt64(&t.isActive) == 1
}

func (t *Transaction) deactivate() {
	atomic.SwapInt64(&t.isActive, 0)
}
