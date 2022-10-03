package sql

import (
	"context"
	"database/sql"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// NewDefaultFactory creates default transaction.Transaction(sqlx.Tx).
func NewDefaultFactory(db *sql.DB) transaction.TrFactory {
	return func() (transaction.Transaction, error) {
		return NewTransaction(context.Background(), nil, db)
	}
}

// Transaction is transaction.Transaction for sql.Tx.
type Transaction struct {
	tr       *sql.Tx
	isActive bool
}

// NewTransaction creates transaction.Transaction for sql.Tx.
func NewTransaction(ctx context.Context, opts *sql.TxOptions, db *sql.DB) (*Transaction, error) {
	tr, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, multierr.Combine(transaction.ErrTransaction, err)
	}

	return &Transaction{tr: tr, isActive: true}, nil
}

// Transaction returns the real transaction sqlx.Tx.
func (t *Transaction) Transaction() interface{} {
	return t.tr
}

// Commit calls close for a database.
func (t *Transaction) Commit() error {
	t.isActive = false

	if err := t.tr.Commit(); err != nil {
		return multierr.Combine(transaction.ErrCommit, err)
	}

	return nil
}

// Rollback calls close for a database.
func (t *Transaction) Rollback() error {
	t.isActive = false

	if err := t.tr.Rollback(); err != nil {
		return multierr.Combine(transaction.ErrRollback, err)
	}

	return nil
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.isActive
}
