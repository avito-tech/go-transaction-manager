package manager

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

// Closer closes trm.Transaction.
type Closer func(context.Context, interface{}, *error) error

type trCloser struct {
	tr     trm.Transaction
	cancel context.CancelFunc
	log    logger
}

func newTxCommit(tr trm.Transaction, l logger, c context.CancelFunc) Closer {
	return (&trCloser{
		tr:     tr,
		cancel: c,
		log:    l,
	}).close
}

//nolint:funlen
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
	isErrSkippable := hasError && trm.IsSkippable(*errInProcessTr)
	// TODO not sure that context errors should be propagated.
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
			if isCtxErr || errors.Is(*errInProcessTr, trm.ErrAlreadyClosed) {
				return *errInProcessTr
			}

			return multierr.Combine(*errInProcessTr, trm.ErrAlreadyClosed)
		}

		return trm.ErrAlreadyClosed
	}

	if hasError && !isErrSkippable {
		if errRollback := c.tr.Rollback(ctx); errRollback != nil {
			return multierr.Combine(*errInProcessTr, trm.ErrRollback, errRollback)
		}

		return *errInProcessTr
	}

	if err := c.tr.Commit(ctx); err != nil {
		var errUnSkipped error
		if isErrSkippable {
			errUnSkipped = trm.UnSkippable(*errInProcessTr)
		}

		return multierr.Combine(trm.ErrCommit, err, errUnSkipped)
	} else if isErrSkippable {
		return *errInProcessTr
	}

	return nil
}

func newNilClose(cancel context.CancelFunc) Closer {
	return func(_ context.Context, p interface{}, err *error) error {
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
