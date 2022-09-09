package transaction

import (
	"context"

	"go.uber.org/multierr"
)

// DefaultCtxKey is a default key to store Transaction.
var DefaultCtxKey = ctxKey{} //nolint:gochecknoglobals

type ctxKey struct{}

// ManagerOpt is a type to configure ManagerImpl.
// TODO is it necessary?
type ManagerOpt func(*ManagerImpl)

// ManagerImpl is an implementation of Manager based on storing Transaction in context.Context.
// TODO rename.
type ManagerImpl struct {
	factory Factory
	key     CtxKey
	log     logger
}

// NewManager creates ManagerImpl.
func NewManager(f Factory) *ManagerImpl {
	return NewManagerOpts(f, NewSettings())
}

// NewManagerOpts creates ManagerImpl with Settings.
func NewManagerOpts(f Factory, settings Settings) *ManagerImpl {
	// TODO implements other settings
	m := &ManagerImpl{
		factory: f,
		key:     settings.CtxKey(),
		log:     defaultLog,
	}

	return m
}

// Do processes a transaction inside a closure.
func (m *ManagerImpl) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	ctx, closer, err := m.init(ctx)
	if err != nil {
		return err
	}
	// Pointer to error is required for recovery and subsequent Transaction.Rollback call.
	defer closer(ctx, &err) //nolint:errcheck // The error will be processed by the caller of Manager.Do.

	return fn(ctx)
}

type closer func(context.Context, *error) error

func (m *ManagerImpl) init(ctx context.Context) (context.Context, closer, error) {
	// TODO add propagation
	tr := TrFromCtx(ctx, m.key)

	if tr == nil {
		tr, err := m.factory()
		if err != nil {
			return nil, nil, multierr.Combine(ErrBegin, err)
		}

		return ctxWithTr(ctx, m.key, tr), newTxCommit(tr, m.log), nil
	}

	return ctx, newNilClose(), nil
}

// WithKey sets CtxKey for ManagerImpl.
func WithKey(key CtxKey) ManagerOpt {
	return func(m *ManagerImpl) {
		m.key = key
	}
}
