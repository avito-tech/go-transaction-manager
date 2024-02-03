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
	tx            txInterface
	active        *drivers.IsClose
	activeClosure *drivers.IsClose
}

// NewTransaction creates trm.Transaction for sqlx.Tx.
func NewTransaction(
	ctx context.Context,
	db redis.UniversalClient,
	s Settings,
) (context.Context, *Transaction, error) {
	t := &Transaction{
		tx:            nil,
		active:        drivers.NewIsClosed(),
		activeClosure: drivers.NewIsClosed(),
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
					tx:        rtx,
					Pipeliner: pipe,
				}

				for _, d := range s.TxDecorators() {
					t.tx = d(t.tx, db)
				}

				wg.Done()

				<-t.activeClosure.Closed()

				return t.activeClosure.Err()
			})

			if len(cmds) > 0 && s.ReturnPtr() != nil {
				s.AppendReturn(cmds...)
			}

			return err
		}, s.WatchKeys()...)

		if t.tx != nil {
			t.active.CloseWithCause(err)
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

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	select {
	case <-t.active.Closed():
		cmds, err := t.tx.Exec(ctx)

		// TODO process cmds
		_ = cmds

		return err
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
		return t.tx.Discard()
	default:
		t.activeClosure.CloseWithCause(drivers.ErrRollbackTr)

		<-t.active.Closed()

		err := t.active.Err()
		if errors.Is(err, drivers.ErrRollbackTr) {
			return nil
		}

		// unreachable code, because of go-redis doesn't process error from Close
		// https://github.com/redis/go-redis/blob/v8.11.5/tx.go#L69
		// https://github.com/redis/go-redis/blob/v8.11.5/pipeline.go#L130

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
