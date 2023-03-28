package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/trm"
)

// Transaction is trm.Transaction for sql.Tx.
// trm.NestedTrFactory returns IsActive as true while trm.Transaction is opened.
type Transaction struct {
	tx        *sql.Tx
	savePoint SavePoint
	saves     int64
	isActive  int64
}

// NewTransaction creates trm.Transaction for sql.Tx.
func NewTransaction(
	ctx context.Context,
	sp SavePoint,
	opts *sql.TxOptions,
	db *sql.DB,
) (context.Context, *Transaction, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{tx: tx, savePoint: sp, isActive: 1, saves: 0}

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

// Transaction returns the real transaction sqlx.Tx.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction by save point.
func (t *Transaction) Begin(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) { //nolint:ireturn,nolintlint
	_, err := t.tx.ExecContext(ctx, t.savePoint.Create(t.incrementID()))
	if err != nil {
		return ctx, nil, err
	}

	return ctx, t, nil
}

// Commit the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	if t.hasSavePoint() {
		_, err := t.tx.ExecContext(ctx, t.savePoint.Release(t.decrementID()))
		if err != nil {
			return multierr.Combine(trm.ErrNestedCommit, err)
		}

		return nil
	}

	defer t.deactivate()

	return t.tx.Commit()
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	if t.hasSavePoint() {
		_, err := t.tx.ExecContext(ctx, t.savePoint.Rollback(t.decrementID()))
		if err != nil {
			return multierr.Combine(trm.ErrNestedRollback, err)
		}

		return nil
	}

	defer t.deactivate()

	return t.tx.Rollback()
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return atomic.LoadInt64(&t.isActive) == 1
}

func (t *Transaction) deactivate() {
	atomic.SwapInt64(&t.isActive, 0)
}

func (t *Transaction) hasSavePoint() bool {
	return atomic.LoadInt64(&t.saves) > 0
}

func (t *Transaction) incrementID() string {
	atomic.AddInt64(&t.saves, 1)

	return t.id()
}

func (t *Transaction) decrementID() string {
	defer atomic.AddInt64(&t.saves, -1)

	return t.id()
}

func (t *Transaction) id() string {
	return fmt.Sprintf("tx_%d", atomic.LoadInt64(&t.saves))
}
