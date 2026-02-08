package manager

import "github.com/avito-tech/go-transaction-manager/trm/v2"

// WithDefaultCommitHooks sets manager-level default commit hooks, registered before InitialHooks.
func WithDefaultCommitHooks(hooks []trm.CommitHook, opts ...trm.HookOption) ManagerXOpt {
	return func(xm *XManager) error {
		xm.defaultCommitHooks = hooks
		xm.defaultCommitOpts = opts
		return nil
	}
}

// WithDefaultRollbackHooks sets manager-level default rollback hooks, registered before InitialHooks.
func WithDefaultRollbackHooks(hooks []trm.RollbackHook, opts ...trm.HookOption) ManagerXOpt {
	return func(xm *XManager) error {
		xm.defaultRollbackHooks = hooks
		xm.defaultRollbackOpts = opts
		return nil
	}
}
