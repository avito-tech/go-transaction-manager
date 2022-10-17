package manager

import (
	"context"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

type trCloser struct {
	tr transaction.Transaction
	// cancel context.CancelFunc
	log logger
}

func newTxCommit(tr transaction.Transaction, l logger) closer {
	return (&trCloser{
		tr:  tr,
		log: l,
	}).close
}

func (c *trCloser) close(_ context.Context, p interface{}, errInProcessTr *error) error {
	// defer c.cancel()
	// recover from panic
	if p != nil {
		if err := c.tr.Rollback(); err != nil {
			c.log.Printf("%v, %v", err, p)
		}

		panic(p)
	}

	if *errInProcessTr != nil {
		if errRollback := c.tr.Rollback(); errRollback != nil {
			return multierr.Combine(*errInProcessTr, transaction.ErrRollback, errRollback)
		}

		return *errInProcessTr
	}

	if err := c.tr.Commit(); err != nil {
		return multierr.Combine(transaction.ErrCommit, err)
	}

	return nil
}

func newNilClose() closer {
	return func(ctx context.Context, p interface{}, err *error) error {
		if p != nil {
			panic(p)
		}

		if *err != nil {
			return *err
		}

		return nil
	}
}
