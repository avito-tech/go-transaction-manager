// Package manager implements a transaction.Manager interface.
package manager

import (
	"context"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/transaction"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

// Opt is a type to configure Manager.
type Opt func(*Manager)

// Manager is an implementation of Manager based on storing Transaction in context.Context.
type Manager struct {
	factory  transaction.TrFactory
	settings transaction.Settings
	log      logger
}

// New creates Manager.
func New(f transaction.TrFactory, oo ...Opt) *Manager {
	m := &Manager{
		factory:  f,
		log:      defaultLog,
		settings: settings.New(),
	}

	for _, o := range oo {
		o(m)
	}

	return m
}

// WithSettings sets transaction.Settings for Manager.
func WithSettings(s transaction.Settings) Opt {
	return func(m *Manager) {
		m.settings = s
	}
}

// Do processes a transaction inside a closure.
func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	ctx, closer, err := m.init(ctx)
	if err != nil {
		return err
	}
	// Pointer to error is required for recovery and subsequent Transaction.Rollback call.
	defer closer(ctx, &err) //nolint:errcheck // The error will be processed by the caller of Manager.Do.

	return fn(ctx)
}

type closer func(context.Context, *error) error

func (m *Manager) init(ctx context.Context) (context.Context, closer, error) {
	tr := transaction.TrFromCtx(ctx, m.settings.CtxKey())
	isOpened := tr == nil

	switch m.settings.Propagation() {
	case transaction.PropagationNever:
		if isOpened {
			return ctx, nil, transaction.ErrPropagationNever
		}

		return ctx, newNilClose(), nil
	case transaction.PropagationsMandatory:
		if isOpened {
			return ctx, nil, transaction.ErrPropagationMandatory
		}
	case transaction.PropagationRequired:
		if isOpened {
			return ctx, newNilClose(), nil
		}
	case transaction.PropagationNotSupported:
		if isOpened {
			// TODO remove transaction from ctx
			panic("todo")
		}

		return ctx, newNilClose(), nil
	case transaction.PropagationRequiresNew:
		if isOpened {
			// TODO remove transaction from ctx
			panic("todo")
		}
	case transaction.PropagationSupports:
		// TODO
		panic("todo")
	case transaction.PropagationNested:
		// TODO create nested transaction
		panic("todo")
	}

	tr, err := m.factory()
	if err != nil {
		return nil, nil, multierr.Combine(transaction.ErrBegin, err)
	}

	return transaction.CtxWithTr(ctx, m.settings.CtxKey(), tr), newTxCommit(tr, m.log), nil
}
