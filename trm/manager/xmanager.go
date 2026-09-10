package manager

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type xSettingsCtxKey struct{}

// TODO: check if it is required.
func withXSettings(ctx context.Context, xs trm.XSettings) context.Context {
	return context.WithValue(ctx, xSettingsCtxKey{}, xs)
}

// TODO: check if it is required.
func XSettingsFromContext(ctx context.Context) trm.XSettings {
	if xs, ok := ctx.Value(xSettingsCtxKey{}).(trm.XSettings); ok {
		return xs
	}

	return nil
}

// ManagerXOpt configures XManager.
type ManagerXOpt func(*XManager) error

// XManager provides XDo and XDoWithSettings with commit/rollback hooks and TxInfo.
// It wraps Manager and does not modify driver transactions.
type XManager struct {
	manager *Manager

	defaultCommitHooks   []trm.CommitHook
	defaultCommitOpts    []trm.HookOption
	defaultRollbackHooks []trm.RollbackHook
	defaultRollbackOpts  []trm.HookOption
}

// NewXManager creates an XManager from an existing Manager.
func NewXManager(m *Manager, oo ...ManagerXOpt) (*XManager, error) {
	xManager := &XManager{
		manager:              m,
		defaultCommitHooks:   nil,
		defaultCommitOpts:    nil,
		defaultRollbackHooks: nil,
		defaultRollbackOpts:  nil,
	}
	for _, o := range oo {
		if err := o(xManager); err != nil {
			return nil, err
		}
	}

	return xManager, nil
}

// MustXManager returns XManager if err is nil and panics otherwise.
func MustXManager(m *Manager, oo ...ManagerXOpt) *XManager {
	xm, err := NewXManager(m, oo...)
	if err != nil {
		panic(err)
	}

	return xm
}

// XDo runs fn inside a transaction with default XSettings (no savepoint hooks, no initial hooks).
func (xm *XManager) XDo(ctx context.Context, fn func(ctx context.Context) error) error {
	return xm.XDoWithSettings(ctx, trm.DefaultXSettings(xm.manager.settings), fn)
}

// XDoWithSettings runs fn inside a transaction with the given XSettings.
// The transaction in context is an XTransaction; RegisterCommit/RegisterRollback can be used inside fn.
func (xm *XManager) XDoWithSettings(
	ctx context.Context,
	xs trm.XSettings,
	fn func(ctx context.Context) error,
) (err error) {
	s := xs.EnrichBy(xm.manager.settings)
	enableSavepointHooks := xs.EnableSavepointHooks()

	factory := func(ctx context.Context, s trm.Settings) (context.Context, trm.Transaction, error) {
		ctx, tr, err := xm.manager.getTransaction(ctx, s)
		if err != nil {
			return ctx, nil, err
		}
		xtx := newXTransaction(tr, s.Propagation(), enableSavepointHooks)

		return ctx, xtx, nil
	}

	ctx = withXSettings(ctx, xs)

	ctx, closer, err := xm.manager.initWithFactory(ctx, s, factory)
	if err != nil {
		return err
	}

	tr := xm.manager.ctxManager.ByKey(ctx, s.CtxKey())
	if x, ok := tr.(*xTransaction); ok {
		ctx = trm.WithXTransaction(ctx, x)
		if x.info.NestingLevel() == 0 {
			registerDefaults(xm, x)
			registerInitialHooks(x, xs.InitialHooks())
		}
	}

	defer func() { err = closer(ctx, recover(), &err) }()

	return fn(ctx)
}

func registerDefaults(xManager *XManager, xTr *xTransaction) {
	opts := trm.ApplyOpts(xManager.defaultCommitOpts)
	if len(xManager.defaultCommitHooks) > 0 {
		xTr.RegisterCommitHooks(xManager.defaultCommitHooks, opts)
	}

	opts = trm.ApplyOpts(xManager.defaultRollbackOpts)
	if len(xManager.defaultRollbackHooks) > 0 {
		xTr.RegisterRollbackHooks(xManager.defaultRollbackHooks, opts)
	}
}

func registerInitialHooks(xTr *xTransaction, initialHooks []trm.Hooks) {
	for _, hook := range initialHooks {
		opts := trm.ApplyOpts(hook.Opts)
		if len(hook.Commits) > 0 {
			xTr.RegisterCommitHooks(hook.Commits, opts)
		}

		if len(hook.Rollbacks) > 0 {
			xTr.RegisterRollbackHooks(hook.Rollbacks, opts)
		}
	}
}
