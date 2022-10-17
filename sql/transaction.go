package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// Transaction is transaction.Transaction for sql.Tx.
// transaction.SavePoint has IsActive as true while transaction.Transaction is opened.
type Transaction struct {
	tx        *sql.Tx
	savePoint transaction.SavePoint
	saves     int64
	isActive  bool
}

// NewTransaction creates transaction.Transaction for sql.Tx.
func NewTransaction(
	ctx context.Context,
	sp transaction.SavePoint,
	opts *sql.TxOptions,
	db *sql.DB,
) (*Transaction, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, multierr.Combine(transaction.ErrTransaction, err)
	}

	return &Transaction{tx: tx, savePoint: sp, isActive: true}, nil
}

// Transaction returns the real transaction sqlx.Tx.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// SavePoint creates nested transaction by save point.
func (t *Transaction) SavePoint(ctx context.Context, _ transaction.Settings) (*Transaction, error) {
	// TODO check that is transaction.Settings necessary
	_, err := t.tx.ExecContext(ctx, t.savePoint.Create(t.incrementID()))
	if err != nil {
		return nil, multierr.Combine(transaction.ErrTransaction, err)
	}

	return t, nil
}

// Commit calls close for a database.
func (t *Transaction) Commit() error {
	t.isActive = false

	if t.hasSavePoint() {
		_, err := t.tx.Exec(t.savePoint.Release(t.id()))
		t.decrementID()

		if err != nil {
			return multierr.Combine(transaction.ErrCommit, err)
		}

		return nil
	}

	if err := t.tx.Commit(); err != nil {
		return multierr.Combine(transaction.ErrCommit, err)
	}

	return nil
}

// Rollback calls close for a database.
func (t *Transaction) Rollback() error {
	if t.hasSavePoint() {
		_, err := t.tx.Exec(t.savePoint.Rollback(t.id()))
		if err != nil {
			return multierr.Combine(transaction.ErrCommit, err)
		}

		return nil
	}

	t.isActive = false

	if err := t.tx.Rollback(); err != nil {
		return multierr.Combine(transaction.ErrRollback, err)
	}

	return nil
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.isActive
}

func (t *Transaction) hasSavePoint() bool {
	return atomic.LoadInt64(&t.saves) > 0
}

func (t *Transaction) incrementID() string {
	atomic.AddInt64(&t.saves, 1)

	return t.id()
}

func (t *Transaction) decrementID() string {
	atomic.AddInt64(&t.saves, -1)

	return t.id()
}

func (t *Transaction) id() string {
	return fmt.Sprintf("tx_%d", atomic.LoadInt64(&t.saves))
}
