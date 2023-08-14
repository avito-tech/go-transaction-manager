package pgx

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jackc/pgx/v4"
	"go.uber.org/multierr"

	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"github.com/avito-tech/go-transaction-manager/trm"
)

// Transaction is trm.Transaction for pgx.Tx.
type Transaction struct {
	tx        pgx.Tx
	savePoint trmsql.SavePoint
	saves     int64
	isActive  int64
}

// NewTransaction creates trm.Transaction for pgx.Tx.
func NewTransaction(
	ctx context.Context,
	sp trmsql.SavePoint,
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
		tx:        tx,
		savePoint: sp,
		isActive:  1,
		saves:     0,
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
	_, err := t.tx.Exec(ctx, t.savePoint.Create(t.incrementID()))
	if err != nil {
		return ctx, nil, err
	}

	return ctx, t, nil
}

// Commit the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	if t.hasSavePoint() {
		_, err := t.tx.Exec(ctx, t.savePoint.Release(t.decrementID()))
		if err != nil {
			return multierr.Combine(trm.ErrNestedCommit, err)
		}

		return nil
	}

	defer t.deactivate()

	return t.tx.Commit(ctx)
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	if t.hasSavePoint() {
		_, err := t.tx.Exec(ctx, t.savePoint.Rollback(t.decrementID()))
		if err != nil {
			return multierr.Combine(trm.ErrNestedRollback, err)
		}

		return nil
	}

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
