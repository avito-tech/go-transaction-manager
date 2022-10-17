package sqlx

// TODO move common solutions for sqlx and sql in one place.

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// Transaction is transaction.Transaction for sqlx.Tx.
type Transaction struct {
	tx        *sqlx.Tx
	savePoint transaction.SavePoint
	saves     int64
	isActive  bool
}

// NewTransaction creates transaction.Transaction for sqlx.Tx.
func NewTransaction(
	ctx context.Context,
	sp transaction.SavePoint,
	opts *sql.TxOptions,
	db *sqlx.DB,
) (*Transaction, error) {
	tx, err := db.BeginTxx(ctx, opts)
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
func (t *Transaction) SavePoint(ctx context.Context, _ transaction.Settings) (transaction.Transaction, error) { //nolint:ireturn
	// TODO check that is transaction.Settings necessary.
	_, err := t.tx.ExecContext(ctx, t.savePoint.Create(t.incrementID()))
	if err != nil {
		return nil, multierr.Combine(transaction.ErrSPBegin, err)
	}

	return t, nil
}

// Commit calls close for a database.
func (t *Transaction) Commit() error {
	if t.hasSavePoint() {
		_, err := t.tx.Exec(t.savePoint.Release(t.decrementID()))
		if err != nil {
			return multierr.Combine(transaction.ErrSPCommit, err)
		}

		return nil
	}

	t.isActive = false

	if err := t.tx.Commit(); err != nil {
		return multierr.Combine(transaction.ErrCommit, err)
	}

	return nil
}

// Rollback calls close for a database.
func (t *Transaction) Rollback() error {
	if t.hasSavePoint() {
		_, err := t.tx.Exec(t.savePoint.Rollback(t.decrementID()))
		if err != nil {
			return multierr.Combine(transaction.ErrSPRollback, err)
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
	defer atomic.AddInt64(&t.saves, -1)

	return t.id()
}

func (t *Transaction) id() string {
	return fmt.Sprintf("tx_%d", atomic.LoadInt64(&t.saves))
}
