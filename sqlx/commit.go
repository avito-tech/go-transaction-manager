package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// Commit completes the transaction.
type Commit func(*error) error

type trCommit struct {
	tr     *sqlx.Tx
	cancel context.CancelFunc
	log    logger
}

func newTxCommit(tr *sqlx.Tx, cancel context.CancelFunc, l logger) Commit {
	return (&trCommit{
		tr:     tr,
		cancel: cancel,
		log:    l,
	}).commit
}

func (c *trCommit) commit(errInProcessTr *error) error {
	defer c.cancel()

	if p := recover(); p != nil {
		if err := c.tr.Rollback(); err != nil {
			c.log.Printf("%v, %v", err, p)
		}

		panic(p)
	}

	if *errInProcessTr != nil {
		if errRollback := c.tr.Rollback(); errRollback != nil {
			return multierr.Combine(*errInProcessTr, transaction.ErrTransaction, errRollback)
		}
	}

	if err := c.tr.Commit(); err != nil {
		return multierr.Combine(transaction.ErrTransaction, err)
	}

	return nil
}

func newNilCommit(cancel context.CancelFunc) Commit {
	return func(err *error) error {
		defer cancel()

		if *err != nil {
			return *err
		}

		return nil
	}
}
