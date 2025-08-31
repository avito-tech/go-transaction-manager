// Package manager implements a trm.Manager interface.
package manager

import (
	"context"

	"go.uber.org/multierr"

	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
)

// Opt is a type to configure Manager.
type Opt func(*Manager) error

// Manager is an implementation of Manager based on storing Transaction in context.Context.
type Manager struct {
	getTransaction trm.TrFactory
	settings       trm.Settings
	ctxManager     trm.CtxManager
	log            logger
}

// New creates Manager.
func New(f trm.TrFactory, oo ...Opt) (*Manager, error) {
	s, err := settings.New()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		getTransaction: f,
		log:            defaultLog,
		ctxManager:     trmcontext.DefaultManager,
		settings:       s,
	}

	for _, o := range oo {
		if err := o(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Must returns Manager if err is nil and panics otherwise.
func Must(f trm.TrFactory, oo ...Opt) *Manager {
	s, err := New(f, oo...)
	if err != nil {
		panic(err)
	}

	return s
}

// Do processes a transaction inside a closure.
func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	return m.DoWithSettings(ctx, m.settings, fn)
}

// DoWithSettings processes a transaction inside a closure with custom trm.Settings.
func (m *Manager) DoWithSettings(ctx context.Context, s trm.Settings, fn func(ctx context.Context) error) (err error) {
	ctx, closer, err := m.Init(ctx, s.EnrichBy(m.settings))
	if err != nil {
		return err
	}

	// Pointer to error is required for recovery and subsequent trm.Transaction rollback call.
	defer func() {
		err = closer(ctx, recover(), &err)
	}()

	return fn(ctx)
}

// Init creates a context.Context with a trm.Transaction and Closer to finish trm.Transaction.
// Required to explicitly close the transaction by calling Closer.
// Nested goroutines would be canceled after the transaction closing by context.CancelFunc.
func (m *Manager) Init(ctx context.Context, s trm.Settings) (context.Context, Closer, error) {
	tr := m.ctxManager.ByKey(ctx, s.CtxKey())
	isOpened := tr != nil

	ctx, cancel := m.withCancel(ctx, s)

	switch s.Propagation() {
	case trm.PropagationRequired:
		if isOpened {
			return ctx, newNilClose(cancel), nil
		}
	case trm.PropagationNested:
		if isOpened {
			return m.propagationNested(ctx, s, tr, cancel)
		}
	case trm.PropagationsMandatory:
		if isOpened {
			return ctx, newNilClose(cancel), nil
		}

		return ctx, nil, trm.ErrPropagationMandatory
	case trm.PropagationNever:
		if isOpened {
			return ctx, nil, trm.ErrPropagationNever
		}

		return ctx, newNilClose(cancel), nil
	case trm.PropagationNotSupported:
		if isOpened {
			return m.ctxManager.SetByKey(ctx, s.CtxKey(), nil),
				newNilClose(cancel),
				nil
		}

		return ctx, newNilClose(cancel), nil
	case trm.PropagationRequiresNew:
		// do nothing
	case trm.PropagationSupports:
		return ctx, newNilClose(cancel), nil
	}

	ctx, tr, err := m.getTransaction(ctx, s)
	if err != nil {
		return nil, nil, multierr.Combine(trm.ErrBegin, err)
	}

	return m.ctxManager.SetByKey(ctx, s.CtxKey(), tr),
		newTxCommit(tr, m.log, cancel),
		nil
}

func (m *Manager) propagationNested(ctx context.Context, s trm.Settings, tr trm.Transaction, c context.CancelFunc) (context.Context, Closer, error) {
	nestedFactory, ok := tr.(trm.NestedTrFactory)
	if ok {
		ctx, tr, err := nestedFactory.Begin(ctx, s)
		if err != nil {
			return ctx, nil, multierr.Combine(trm.ErrNestedBegin, err)
		}

		return m.ctxManager.SetByKey(ctx, s.CtxKey(), tr),
			newTxCommit(tr, m.log, c),
			nil
	}

	return ctx, newNilClose(c), nil
}

func (m *Manager) withCancel(ctx context.Context, s trm.Settings) (context.Context, context.CancelFunc) {
	t := s.TimeoutOrNil()
	if t != nil {
		return context.WithTimeout(ctx, *t)
	}

	if s.Cancelable() {
		return context.WithCancel(ctx)
	}

	return ctx, nilCancel
}

func nilCancel() {}
