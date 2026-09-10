package manager

import (
	"context"
	"sync/atomic"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

// txInfoImpl is a thread-safe TxInfo.
type txInfoImpl struct {
	propagation trm.Propagation
	level       int32
}

func (t *txInfoImpl) Propagation() trm.Propagation { return t.propagation }
func (t *txInfoImpl) IsNew() bool                  { return atomic.LoadInt32(&t.level) == 0 }
func (t *txInfoImpl) IsNested() bool               { return atomic.LoadInt32(&t.level) > 0 }
func (t *txInfoImpl) NestingLevel() int            { return int(atomic.LoadInt32(&t.level)) }

func (t *txInfoImpl) increment() { atomic.AddInt32(&t.level, 1) }
func (t *txInfoImpl) decrement() { atomic.AddInt32(&t.level, -1) }

// xTransaction wraps trm.Transaction and adds TxInfo and hook registries.
type xTransaction struct {
	underlying           trm.Transaction
	info                 *txInfoImpl
	registry             *trm.HookRegistry
	enableSavepointHooks bool
}

func newXTransaction(
	underlying trm.Transaction,
	propagation trm.Propagation,
	enableSavepointHooks bool,
) *xTransaction {
	return &xTransaction{
		underlying:           underlying,
		info:                 &txInfoImpl{propagation: propagation, level: 0},
		registry:             trm.NewHookRegistry(),
		enableSavepointHooks: enableSavepointHooks,
	}
}

func (x *xTransaction) Transaction() interface{} { return x.underlying.Transaction() }

//nolint:ireturn
func (x *xTransaction) TxInfo() trm.TxInfo { return x.info }

func (x *xTransaction) IsActive() bool { return x.underlying.IsActive() }

func (x *xTransaction) Closed() <-chan struct{} { return x.underlying.Closed() }

func (x *xTransaction) Commit(ctx context.Context) error {
	level := x.info.NestingLevel()
	if level > 0 {
		if x.enableSavepointHooks {
			if err := x.registry.RunSavepointCommitHooks(ctx, x); err != nil {
				_ = x.underlying.Rollback(ctx)

				return err
			}

			x.registry.PopSavepoint()
		} else {
			x.registry.PopSavepoint()
		}

		x.info.decrement()

		return x.underlying.Commit(ctx)
	}

	if err := x.registry.RunTransactionCommitHooks(ctx, x); err != nil {
		_ = x.underlying.Rollback(ctx)

		return err
	}

	return x.underlying.Commit(ctx)
}

func (x *xTransaction) Rollback(ctx context.Context) error {
	level := x.info.NestingLevel()
	if level > 0 {
		if x.enableSavepointHooks {
			x.registry.RunSavepointRollbackHooks(ctx, x)
			x.registry.PopSavepoint()
		} else {
			x.registry.PopSavepoint()
		}

		x.info.decrement()

		return x.underlying.Rollback(ctx)
	}

	x.registry.RunTransactionRollbackHooks(ctx, x)

	return x.underlying.Rollback(ctx)
}

func (x *xTransaction) Begin(ctx context.Context, s trm.Settings) (
	context.Context,
	trm.Transaction,
	error,
) {
	nested, ok := x.underlying.(trm.NestedTrFactory)
	if !ok {
		return ctx, x, nil
	}

	ctx, _, err := nested.Begin(ctx, s)
	if err != nil {
		return ctx, nil, err
	}

	x.registry.PushSavepoint()
	x.info.increment()

	return ctx, x, nil
}

func (x *xTransaction) RegisterCommitHooks(hooks []trm.CommitHook, opts *trm.HookOpts) {
	x.registry.AddCommitHooks(hooks, opts)
}

func (x *xTransaction) RegisterRollbackHooks(hooks []trm.RollbackHook, opts *trm.HookOpts) {
	x.registry.AddRollbackHooks(hooks, opts)
}
