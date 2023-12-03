//go:build go1.16
// +build go1.16

// Package gorm is an implementation of trm.Transaction interface by Transaction for *gorm.DB.
package gorm

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"sync/atomic"

	"gorm.io/gorm"

	trm "github.com/avito-tech/go-transaction-manager/v2"
)

var errRollbackTx = errors.New("rollback transaction")

// Transaction is trm.Transaction for sqlx.Tx.
type Transaction struct {
	tx       *gorm.DB
	err      chan error
	isActive int64
}

// NewTransaction creates trm.Transaction for sqlx.Tx.
func NewTransaction(
	ctx context.Context,
	opts *sql.TxOptions,
	db *gorm.DB,
) (context.Context, *Transaction, error) {
	tr := &Transaction{isActive: 1, err: make(chan error), tx: nil}

	var err error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		db = db.WithContext(ctx)
		err = db.Transaction(func(tx *gorm.DB) error {
			tr.tx = tx

			wg.Done()

			return <-tr.err
		}, opts)

		if tr.tx != nil {
			tr.err <- err
		} else {
			wg.Done()
		}
	}()

	wg.Wait()

	if err != nil {
		return ctx, nil, err
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

// Transaction returns the real transaction sqlx.Tx.
// trm.NestedTrFactory returns IsActive as true while trm.Transaction is opened.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction by save point.
func (t *Transaction) Begin(ctx context.Context, s trm.Settings) (context.Context, trm.Transaction, error) {
	return NewDefaultFactory(t.tx)(ctx, s)
}

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(_ context.Context) error {
	defer t.deactivate()

	t.err <- nil

	return <-t.err
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(_ context.Context) error {
	defer t.deactivate()

	t.err <- errRollbackTx

	err := <-t.err

	if errors.Is(err, errRollbackTx) {
		return nil
	}

	return err
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return atomic.LoadInt64(&t.isActive) == 1
}

func (t *Transaction) deactivate() {
	atomic.SwapInt64(&t.isActive, 0)
}
