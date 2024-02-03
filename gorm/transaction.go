//go:build go1.16
// +build go1.16

// Package gorm is an implementation of trm.Transaction interface by Transaction for *gorm.DB.
package gorm

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"gorm.io/gorm"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/drivers"
)

// Transaction is trm.Transaction for sqlx.Tx.
type Transaction struct {
	tx            *gorm.DB
	txMutex       sync.Mutex
	active        *drivers.IsClose
	activeClosure *drivers.IsClose
}

// NewTransaction creates trm.Transaction for sqlx.Tx.
func NewTransaction(
	ctx context.Context,
	opts *sql.TxOptions,
	db *gorm.DB,
) (context.Context, *Transaction, error) {
	tr := &Transaction{
		tx:            nil,
		txMutex:       sync.Mutex{},
		active:        drivers.NewIsClosed(),
		activeClosure: drivers.NewIsClosed(),
	}

	var err error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		db = db.WithContext(ctx)
		// Used closure to avoid implementing nested transactions.
		err = db.Transaction(func(tx *gorm.DB) error {
			tr.tx = tx

			wg.Done()

			<-tr.activeClosure.Closed()

			return tr.activeClosure.Err()
		}, opts)

		tr.txMutex.Lock()
		defer tr.txMutex.Unlock()
		tx := tr.tx

		if tx != nil {
			// Return error from transaction rollback
			// Error from commit returns from db.Transaction closure
			if errors.Is(err, drivers.ErrRollbackTr) &&
				tx.Error != nil {
				err = tr.tx.Error
			}

			tr.active.CloseWithCause(err)
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

	select {
	case <-ctx.Done():
		// Rollback will be called by context.Err()
		t.activeClosure.Close()
	case <-t.active.Closed():
	}
}

// Transaction returns the real transaction sqlx.Tx.
// trm.NestedTrFactory returns IsActive as true while trm.Transaction is opened.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction by save point.
func (t *Transaction) Begin(ctx context.Context, s trm.Settings) (context.Context, trm.Transaction, error) {
	t.txMutex.Lock()
	defer t.txMutex.Unlock()

	return NewDefaultFactory(t.tx)(ctx, s)
}

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(_ context.Context) error {
	select {
	case <-t.active.Closed():
		t.txMutex.Lock()
		defer t.txMutex.Unlock()

		return t.tx.Commit().Error
	default:
		t.activeClosure.Close()

		<-t.active.Closed()

		return t.active.Err()
	}
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(_ context.Context) error {
	select {
	case <-t.active.Closed():
		t.txMutex.Lock()
		defer t.txMutex.Unlock()

		return t.tx.Rollback().Error
	default:
		t.activeClosure.CloseWithCause(drivers.ErrRollbackTr)

		<-t.active.Closed()

		err := t.active.Err()
		if errors.Is(err, drivers.ErrRollbackTr) {
			return nil
		}

		return err
	}
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.active.IsActive()
}

func (t *Transaction) Closed() <-chan struct{} {
	return t.active.Closed()
}
