package manager

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type xSettingsCtxKey struct{}

func withXSettings(ctx context.Context, xs trm.XSettings) context.Context {
	return context.WithValue(ctx, xSettingsCtxKey{}, xs)
}

func xSettingsFromContext(ctx context.Context) trm.XSettings {
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
	xm := &XManager{manager: m}
	for _, o := range oo {
		if err := o(xm); err != nil {
			return nil, err
		}
	}
	return xm, nil
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
func (xm *XManager) XDoWithSettings(ctx context.Context, xs trm.XSettings, fn func(ctx context.Context) error) (err error) {
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

func registerDefaults(xm *XManager, x *xTransaction) {
	opts := trm.ApplyOpts(xm.defaultCommitOpts)
	if len(xm.defaultCommitHooks) > 0 {
		x.RegisterCommitHooks(xm.defaultCommitHooks, opts)
	}
	opts = trm.ApplyOpts(xm.defaultRollbackOpts)
	if len(xm.defaultRollbackHooks) > 0 {
		x.RegisterRollbackHooks(xm.defaultRollbackHooks, opts)
	}
}

func registerInitialHooks(x *xTransaction, initial []trm.Hooks) {
	for _, h := range initial {
		opts := trm.ApplyOpts(h.Opts)
		if len(h.Commits) > 0 {
			x.RegisterCommitHooks(h.Commits, opts)
		}
		if len(h.Rollbacks) > 0 {
			x.RegisterRollbackHooks(h.Rollbacks, opts)
		}
	}
}
