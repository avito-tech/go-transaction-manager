package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// Transaction is transaction.Transaction for sqlx.Tx.
type Transaction struct {
	tr       *sqlx.Tx
	isActive bool
}

// NewTransaction creates transaction.Transaction for sqlx.Tx.
func NewTransaction(ctx context.Context, opts *sql.TxOptions, db *sqlx.DB) (*Transaction, error) {
	tr, err := db.BeginTxx(ctx, opts)
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
