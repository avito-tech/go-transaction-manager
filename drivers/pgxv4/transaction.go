package pgxv4

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers"
)

// Transaction is trm.Transaction for pgx.Tx.
//
// Transaction is NOT safe for concurrent use. pgx.Tx does not support running
// commands from multiple goroutines simultaneously (jackc/pgx#2332), so a query
// must never overlap with Commit or Rollback on the same Transaction. For this
// reason context cancellation does not roll back the Transaction from a
// background goroutine: with manager.Manager the rollback is issued after the
// transactional function returns, and standalone callers must call Rollback
// themselves.
type Transaction struct {
	tx       pgx.Tx
	isClosed *drivers.IsClosed
}

func newDefaultTransaction(tx pgx.Tx) *Transaction {
	return &Transaction{
		tx:       tx,
		isClosed: drivers.NewIsClosed(),
	}
}

// NewTransaction creates trm.Transaction for pgx.Tx.
func NewTransaction(
	ctx context.Context,
	opts pgx.TxOptions,
	db Transactional,
) (context.Context, *Transaction, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}

	tr := newDefaultTransaction(tx)

	return ctx, tr, nil
}

// Transaction returns the real transaction pgx.Tx.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction by save point.
func (t *Transaction) Begin(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
	tx, err := t.tx.Begin(ctx)
	if err != nil {
		return ctx, nil, err
	}

	tr := newDefaultTransaction(tx)

	return ctx, tr, nil
}

// Commit the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	defer t.isClosed.Close()

	return t.tx.Commit(ctx)
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	defer t.isClosed.Close()

	return t.tx.Rollback(ctx)
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.isClosed.IsActive()
}

// Closed returns a channel that's closed when transaction committed or rolled back.
func (t *Transaction) Closed() <-chan struct{} {
	return t.isClosed.Closed()
}
