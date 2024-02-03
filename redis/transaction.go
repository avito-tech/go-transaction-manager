// Package redis is an implementation of trm.Transaction interface by Transaction for redis.UniversalClient.
package redis

import (
	"context"
	"errors"
	"sync"

	"github.com/go-redis/redis/v8"

	"github.com/avito-tech/go-transaction-manager/trm/drivers"
)

var errRollbackTx = errors.New("rollback transaction")

// TxDecorator is an interface for Transaction.tx decoration.
type TxDecorator func(tx Cmdable, db redis.Cmdable) Cmdable

// Transaction is trm.Transaction for sqlx.Tx.
type Transaction struct {
	tx Cmdable
	// err is used to close transaction and get error from it
	err    chan error
	active *drivers.IsClose
}

// NewTransaction creates trm.Transaction for sqlx.Tx.
func NewTransaction(
	ctx context.Context,
	db redis.UniversalClient,
	s Settings,
) (context.Context, *Transaction, error) {
	t := &Transaction{
		err:    make(chan error),
		tx:     nil,
		active: drivers.NewIsClosed(),
	}

	var err error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		var cmds []redis.Cmder

		err = db.Watch(ctx, func(rtx *redis.Tx) error {
			fn := rtx.Pipelined
			if s.IsMulti() {
				fn = rtx.TxPipelined
			}

			cmds, err = fn(ctx, func(pipe redis.Pipeliner) error {
				t.tx = &tx{
					tx:      rtx,
					Cmdable: pipe,
				}

				for _, d := range s.TxDecorators() {
					t.tx = d(t.tx, db)
				}

				wg.Done()

				return <-t.err
			})

			if len(cmds) > 0 && s.Return() != nil {
				*s.Return() = append(*s.Return(), cmds...)
			}

			return err
		}, s.WatchKeys()...)

		if t.tx != nil {
			t.err <- err
		} else {
			wg.Done()
		}
	}()

	wg.Wait()

	if err != nil {
		return ctx, nil, err
	}

	go t.awaitDone(ctx)

	return ctx, t, nil
}

func (t *Transaction) awaitDone(ctx context.Context) {
	if ctx.Done() == nil {
		return
	}

	select {
	case <-ctx.Done():
		t.active.Close()
	case <-t.active.Closed():
	}
}

// Transaction returns the real transaction sqlx.Tx.
// trm.NestedTrFactory returns IsActive as true while trm.Transaction is opened.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(_ context.Context) error {
	defer t.active.Close()

	// TODO deadlock
	t.err <- nil

	return <-t.err
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(_ context.Context) error {
	defer t.active.Close()

	// TODO deadlock
	t.err <- errRollbackTx

	err := <-t.err

	if errors.Is(err, errRollbackTx) {
		return nil
	}

	return err
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.active.IsActive()
}
