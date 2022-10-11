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
type Transaction struct {
	tx        *sql.Tx
	savePoint transaction.SavePoint
	saves     int64
	isActive  bool
}

// NewTransaction creates transaction.Transaction for sql.Tx.
func NewTransaction(ctx context.Context, sp transaction.SavePoint, opts *sql.TxOptions, db *sql.DB) (*Transaction, error) {
	tr, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, multierr.Combine(transaction.ErrTransaction, err)
	}

	return &Transaction{tx: tr, savePoint: sp, isActive: true}, nil
}

// Transaction returns the real transaction sqlx.Tx.
func (tr *Transaction) Transaction() interface{} {
	return tr.tx
}

// SavePoint creates nested transaction by save point.
func (tr *Transaction) SavePoint(ctx context.Context, _ transaction.Settings) (*Transaction, error) {
	// TODO check that is transaction.Settings necessary
	_, err := tr.tx.ExecContext(ctx, tr.savePoint.Create(tr.incrementID()))
	if err != nil {
		return nil, multierr.Combine(transaction.ErrTransaction, err)
	}

	return tr, nil
}

// Commit calls close for a database.
func (tr *Transaction) Commit() error {
	tr.isActive = false

	if tr.hasSavePoint() {
		_, err := tr.tx.Exec(tr.savePoint.Release(tr.id()))
		tr.decrementID()

		if err != nil {
			return multierr.Combine(transaction.ErrCommit, err)
		}

		return nil
	}

	if err := tr.tx.Commit(); err != nil {
		return multierr.Combine(transaction.ErrCommit, err)
	}

	return nil
}

// Rollback calls close for a database.
func (tr *Transaction) Rollback() error {
	tr.isActive = false

	if tr.hasSavePoint() {
		_, err := tr.tx.Exec(tr.savePoint.Rollback(tr.id()))
		if err != nil {
			return multierr.Combine(transaction.ErrCommit, err)
		}

		return nil
	}

	if err := tr.tx.Rollback(); err != nil {
		return multierr.Combine(transaction.ErrRollback, err)
	}

	return nil
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (tr *Transaction) IsActive() bool {
	return tr.isActive
}

func (tr *Transaction) hasSavePoint() bool {
	return atomic.LoadInt64(&tr.saves) > 0
}

func (tr *Transaction) incrementID() string {
	atomic.AddInt64(&tr.saves, 1)

	return tr.id()
}

func (tr *Transaction) decrementID() string {
	atomic.AddInt64(&tr.saves, -1)

	return tr.id()
}

func (tr *Transaction) id() string {
	return fmt.Sprintf("tx_%d", atomic.LoadInt64(&tr.saves))
}
