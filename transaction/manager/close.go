package manager

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

type trCloser struct {
	tr     transaction.Transaction
	cancel context.CancelFunc
	log    logger
}

func newTxCommit(tr transaction.Transaction, l logger, c context.CancelFunc) closer {
	return (&trCloser{
		tr:     tr,
		cancel: c,
		log:    l,
	}).close
}

func (c *trCloser) close(ctx context.Context, p interface{}, errInProcessTr *error) error {
	defer c.cancel()

	// recovering from panic
	if p != nil {
		if c.tr.IsActive() {
			if err := c.tr.Rollback(ctx); err != nil {
				c.log.Warning(ctx, fmt.Sprintf("%v, %v", err, p))
			}
		}

		panic(p)
	}

	hasError := *errInProcessTr != nil
	// TODO Not sure that context Errors should be propagated.
	isCtxCanceled := errors.Is(*errInProcessTr, context.Canceled)
	isCtxDeadlineExceeded := errors.Is(*errInProcessTr, context.DeadlineExceeded)
	isCtxErr := isCtxCanceled || isCtxDeadlineExceeded

	ctxErr := ctx.Err()

	if ctxErr != nil {
		if !hasError {
			*errInProcessTr = ctxErr
		} else if !isCtxCanceled && errors.Is(ctxErr, context.Canceled) ||
			!isCtxDeadlineExceeded && errors.Is(ctxErr, context.DeadlineExceeded) {
			*errInProcessTr = multierr.Combine(*errInProcessTr, ctxErr)
		}

		isCtxErr = true
		hasError = true
	}

	if !c.tr.IsActive() {
		if hasError {
			if isCtxErr || errors.Is(*errInProcessTr, transaction.ErrAlreadyClosed) {
				return *errInProcessTr
			}

			return multierr.Combine(*errInProcessTr, transaction.ErrAlreadyClosed)
		}

		return transaction.ErrAlreadyClosed
	}

	if hasError {
		if errRollback := c.tr.Rollback(ctx); errRollback != nil {
			return multierr.Combine(*errInProcessTr, transaction.ErrRollback, errRollback)
		}

		return *errInProcessTr
	}

	if err := c.tr.Commit(ctx); err != nil {
		return multierr.Combine(transaction.ErrCommit, err)
	}

	return nil
}

func newNilClose(cancel context.CancelFunc) closer {
	return func(ctx context.Context, p interface{}, err *error) error {
		defer cancel()

		if p != nil {
			panic(p)
		}

		if *err != nil {
			return *err
		}

		return nil
	}
}
