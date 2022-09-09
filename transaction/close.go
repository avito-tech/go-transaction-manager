package transaction

import (
	"context"

	"go.uber.org/multierr"
)

type trCloser struct {
	tr Transaction
	// cancel context.CancelFunc
	log logger
}

func newTxCommit(tr Transaction, l logger) closer {
	return (&trCloser{
		tr:  tr,
		log: l,
	}).close
}

func (c *trCloser) close(ctx context.Context, errInProcessTr *error) error {
	// defer c.cancel()
	if p := recover(); p != nil {
		if err := c.tr.Rollback(); err != nil {
			c.log.Printf("%v, %v", err, p)
		}

		panic(p)
	}

	if *errInProcessTr != nil {
		if errRollback := c.tr.Rollback(); errRollback != nil {
			return multierr.Combine(*errInProcessTr, ErrRollback, errRollback)
		}
	}

	if err := c.tr.Commit(); err != nil {
		return multierr.Combine(ErrCommit, err)
	}

	return nil
}

func newNilClose() closer {
	return func(ctx context.Context, err *error) error {
		if *err != nil {
			return *err
		}

		return nil
	}
}
