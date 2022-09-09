// Package sqlx is an implementation of transaction.Manager interface.
package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// TrManager is the implementation of transaction.Manager for sqlx.DB.
type TrManager struct {
	db  *sqlx.DB
	log logger
}

// Option is a type to configure TrManager.
type Option func(m *TrManager)

// NewTransactionManager is a factory for TrManager.
func NewTransactionManager(db *sqlx.DB, oo ...Option) *TrManager {
	if db == nil {
		panic("db should be set")
	}

	m := &TrManager{db: db}

	for _, o := range oo {
		o(m)
	}

	return m
}

// Do starts a transaction inside a closure.
func (t *TrManager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	ctx, commit, err := t.Init(ctx)
	if err != nil {
		return err
	}
	// Required for recovery and subsequent call to sqlx.Tx.Rollback
	defer commit(&err) //nolint:errcheck // The error will be processed by caller.

	return fn(ctx)
}

// Init creates a context.Context with a transaction.
// Required to explicitly close the transaction by calling Commit.
// Nested goroutines would be canceled after the transaction closing by context.CancelFunc.
func (t *TrManager) Init(ctx context.Context) (context.Context, Commit, error) {
	ctx, cancel := context.WithCancel(ctx)

	if tr := TrFromCtx(ctx, nil); tr == nil {
		tx, err := t.db.Beginx()
		if err != nil {
			defer cancel()
			return nil, nil, err
		}

		return ctxWithTr(ctx, tx), newTxCommit(tx, cancel, t.log), nil
	}

	return ctx, newNilCommit(cancel), nil
}
