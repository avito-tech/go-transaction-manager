// Package manager implements a transaction.Manager interface.
package manager

import (
	"context"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
	trmcontext "github.com/avito-tech/go-transaction-manager/transaction/context"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

// Opt is a type to configure Manager.
type Opt func(*Manager)

// Manager is an implementation of Manager based on storing Transaction in context.Context.
type Manager struct {
	getTransaction transaction.TrFactory
	settings       transaction.Settings
	ctxManager     transaction.СtxManager
	log            logger
}

// New creates Manager.
func New(f transaction.TrFactory, oo ...Opt) *Manager {
	m := &Manager{
		getTransaction: f,
		log:            defaultLog,
		ctxManager:     trmcontext.DefaultManager,
		settings:       settings.New(),
	}

	for _, o := range oo {
		o(m)
	}

	return m
}

// Do processes a transaction inside a closure.
func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	return m.DoWithSettings(ctx, m.settings, fn)
}

// DoWithSettings processes a transaction inside a closure with custom transaction.Settings.
func (m *Manager) DoWithSettings(ctx context.Context, s transaction.Settings, fn func(ctx context.Context) error) (err error) {
	ctx, closer, err := m.init(ctx, s)
	if err != nil {
		return err
	}

	// Pointer to error is required for recovery and subsequent Transaction.Rollback call.
	defer func() { err = closer(ctx, recover(), &err) }()

	return fn(ctx)
}

type closer func(context.Context, interface{}, *error) error

func (m *Manager) init(ctx context.Context, s transaction.Settings) (context.Context, closer, error) {
	tr := m.ctxManager.ByKey(ctx, s.CtxKey())
	isOpened := tr != nil

	switch s.Propagation() {
	case transaction.PropagationRequired:
		if isOpened {
			return ctx, newNilClose(), nil
		}
	case transaction.PropagationNested:
		if isOpened {
			spFactor, ok := tr.(transaction.SPFactory)
			if ok {
				tr, err := spFactor.SavePoint(ctx, s)
				if err != nil {
					return ctx, nil, err
				}

				return m.ctxManager.SetByKey(ctx, s.CtxKey(), tr),
					newTxCommit(tr, m.log),
					nil
			}

			return ctx, newNilClose(), nil
		}
	case transaction.PropagationsMandatory:
		if isOpened {
			return ctx, newNilClose(), nil
		}

		return ctx, nil, transaction.ErrPropagationMandatory
	case transaction.PropagationNever:
		if isOpened {
			return ctx, nil, transaction.ErrPropagationNever
		}

		return ctx, newNilClose(), nil
	case transaction.PropagationNotSupported:
		if isOpened {
			return m.ctxManager.SetByKey(ctx, s.CtxKey(), nil),
				newNilClose(),
				nil
		}

		return ctx, newNilClose(), nil
	case transaction.PropagationRequiresNew:
		// do nothing
	case transaction.PropagationSupports:
		return ctx, newNilClose(), nil
	}

	tr, err := m.getTransaction(ctx)
	if err != nil {
		return nil, nil, multierr.Combine(transaction.ErrBegin, err)
	}

	return m.ctxManager.SetByKey(ctx, s.CtxKey(), tr),
		newTxCommit(tr, m.log),
		nil
}
