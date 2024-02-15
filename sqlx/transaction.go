package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"

	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/drivers"
)

// Transaction is trm.Transaction for sqlx.Tx.
type Transaction struct {
	tx        *sqlx.Tx
	savePoint trmsql.SavePoint
	saves     int64
	isClosed  *drivers.IsClosed
}

// NewTransaction creates trm.Transaction for sqlx.Tx.
func NewTransaction(
	ctx context.Context,
	sp trmsql.SavePoint,
	opts *sql.TxOptions,
	db *sqlx.DB,
) (context.Context, *Transaction, error) {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{
		tx:        tx,
		savePoint: sp,
		saves:     0,
		isClosed:  drivers.NewIsClosed(),
	}

	go tr.awaitDone(ctx)

	return ctx, tr, nil
}

func (t *Transaction) awaitDone(ctx context.Context) {
	if ctx.Done() == nil {
		return
	}

	select {
	case <-ctx.Done():
		t.isClosed.Close()
	case <-t.isClosed.Closed():
	}
}

// Transaction returns the real transaction sqlx.Tx.
// trm.NestedTrFactory returns IsActive as true while trm.Transaction is opened.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction by save point.
func (t *Transaction) Begin(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
	_, err := t.tx.ExecContext(ctx, t.savePoint.Create(t.incrementID()))
	if err != nil {
		// decrement save point ID after error
		t.decrementID()

		return ctx, nil, multierr.Combine(trm.ErrNestedBegin, err)
	}

	return ctx, t, nil
}

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	if t.hasSavePoint() {
		_, err := t.tx.ExecContext(ctx, t.savePoint.Release(t.decrementID()))
		if err != nil {
			return multierr.Combine(trm.ErrNestedCommit, err)
		}

		return nil
	}

	defer t.isClosed.Close()

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

	defer t.isClosed.Close()

	return t.tx.Rollback()
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.isClosed.IsActive()
}

// Closed returns a channel that's closed when transaction committed or rolled back.
func (t *Transaction) Closed() <-chan struct{} {
	return t.isClosed.Closed()
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
