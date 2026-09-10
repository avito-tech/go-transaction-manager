package trm

import (
	"context"
	"errors"
	"sort"
	"sync"
)

// ErrNoActiveTransaction is returned when no XTransaction is present in context.
var ErrNoActiveTransaction = errors.New("no active XTransaction")

// CommitHook is called before commit. If it returns an error, commit is aborted and rollback is performed.
type CommitHook func(ctx context.Context, tx XTransaction) error

// RollbackHook is called on rollback (including on panic before rethrow).
type RollbackHook func(ctx context.Context, tx XTransaction)

// HookScope defines when hooks run.
type HookScope int

const (
	ScopeTransaction HookScope = iota // run only at real transaction owner completion
	ScopeSavepoint                    // run only for the current nested boundary
)

// HookOption configures hook registration.
type HookOption interface {
	apply(o *HookOpts)
}

// HookOpts holds resolved options for a hook registration (used by XTransaction implementations).
type HookOpts struct {
	Scope HookScope
	Order int
}

type hookOptFunc func(*HookOpts)

func (f hookOptFunc) apply(o *HookOpts) { f(o) }

// ApplyOpts resolves options into HookOpts. Nil or empty opts return default (ScopeTransaction, Order=0).
func ApplyOpts(opts []HookOption) *HookOpts {
	oo := &HookOpts{Scope: ScopeTransaction, Order: 0}
	if opts == nil {
		return oo
	}

	for _, opt := range opts {
		opt.apply(oo)
	}

	return oo
}

func applyOpts(opts []HookOption) *HookOpts {
	return ApplyOpts(opts)
}

// WithScope sets the scope for registered hooks.
//
//nolint:ireturn // returns HookOption interface intentionally
func WithScope(scope HookScope) HookOption {
	return hookOptFunc(func(o *HookOpts) {
		o.Scope = scope
	})
}

// WithOrder sets the execution order; lower runs first. Default is 0. FIFO preserved within same order.
//
//nolint:ireturn // returns HookOption interface intentionally
func WithOrder(order int) HookOption {
	return hookOptFunc(func(o *HookOpts) {
		o.Order = order
	})
}

// Hooks groups commit and rollback hooks with shared options (e.g. for InitialHooks).
type Hooks struct {
	Commits   []CommitHook
	Rollbacks []RollbackHook
	Opts      []HookOption
}

type commitEntry struct {
	hook  CommitHook
	order int
	index int
}

type rollbackEntry struct {
	hook  RollbackHook
	order int
	index int
}

// HookRegistry holds transaction-scope and savepoint-scope hooks.
type HookRegistry struct {
	mu sync.RWMutex

	txCommit   []commitEntry
	txRollback []rollbackEntry

	savepointStack []savepointReg
}

type savepointReg struct {
	commit   []commitEntry
	rollback []rollbackEntry
}

// NewHookRegistry creates a new hook registry.
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		mu:             sync.RWMutex{},
		txCommit:       nil,
		txRollback:     nil,
		savepointStack: nil,
	}
}

func (r *HookRegistry) registerCommit(hooks []CommitHook, opts *HookOpts) {
	if opts == nil {
		opts = &HookOpts{Scope: ScopeTransaction, Order: 0}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	base := len(r.txCommit)
	if opts.Scope == ScopeSavepoint {
		if len(r.savepointStack) == 0 {
			return
		}

		sp := &r.savepointStack[len(r.savepointStack)-1]
		for i, h := range hooks {
			sp.commit = append(sp.commit, commitEntry{hook: h, order: opts.Order, index: base + i})
		}

		return
	}

	for i, h := range hooks {
		r.txCommit = append(r.txCommit, commitEntry{hook: h, order: opts.Order, index: base + i})
	}
}

// AddCommitHooks registers commit hooks on the registry (used by XTransaction implementations).
func (r *HookRegistry) AddCommitHooks(hooks []CommitHook, opts *HookOpts) {
	r.registerCommit(hooks, opts)
}

// AddRollbackHooks registers rollback hooks on the registry.
func (r *HookRegistry) AddRollbackHooks(hooks []RollbackHook, opts *HookOpts) {
	r.registerRollback(hooks, opts)
}

func (r *HookRegistry) registerRollback(hooks []RollbackHook, opts *HookOpts) {
	if opts == nil {
		opts = &HookOpts{Scope: ScopeTransaction, Order: 0}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	base := len(r.txRollback)
	if opts.Scope == ScopeSavepoint {
		if len(r.savepointStack) == 0 {
			return
		}

		sp := &r.savepointStack[len(r.savepointStack)-1]
		for i, h := range hooks {
			sp.rollback = append(
				sp.rollback,
				rollbackEntry{hook: h, order: opts.Order, index: base + i},
			)
		}

		return
	}

	for i, h := range hooks {
		r.txRollback = append(
			r.txRollback,
			rollbackEntry{hook: h, order: opts.Order, index: base + i},
		)
	}
}

// PushSavepoint pushes a new savepoint-level registry for nested transactions.
func (r *HookRegistry) PushSavepoint() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.savepointStack = append(r.savepointStack, savepointReg{
		commit:   nil,
		rollback: nil,
	})
}

// PopSavepoint pops the current savepoint-level registry.
func (r *HookRegistry) PopSavepoint() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.savepointStack) > 0 {
		r.savepointStack = r.savepointStack[:len(r.savepointStack)-1]
	}
}

// RunTransactionCommitHooks runs transaction-scope commit hooks in order. Used by XTransaction implementations.
func (r *HookRegistry) RunTransactionCommitHooks(ctx context.Context, xtx XTransaction) error {
	r.mu.RLock()
	entries := make([]commitEntry, len(r.txCommit))
	copy(entries, r.txCommit)
	r.mu.RUnlock()
	sortCommitEntries(entries)

	for _, e := range entries {
		if err := e.hook(ctx, xtx); err != nil {
			return err
		}
	}

	return nil
}

// RunTransactionRollbackHooks runs transaction-scope rollback hooks in order.
func (r *HookRegistry) RunTransactionRollbackHooks(ctx context.Context, xtx XTransaction) {
	r.mu.RLock()
	entries := make([]rollbackEntry, len(r.txRollback))
	copy(entries, r.txRollback)
	r.mu.RUnlock()
	sortRollbackEntries(entries)

	for _, e := range entries {
		e.hook(ctx, xtx)
	}
}

// RunSavepointCommitHooks runs savepoint-scope commit hooks for the current level. Caller must pop after.
func (r *HookRegistry) RunSavepointCommitHooks(ctx context.Context, xtx XTransaction) error {
	r.mu.RLock()

	var entries []commitEntry
	if len(r.savepointStack) > 0 {
		entries = make([]commitEntry, len(r.savepointStack[len(r.savepointStack)-1].commit))
		copy(entries, r.savepointStack[len(r.savepointStack)-1].commit)
	}

	r.mu.RUnlock()
	sortCommitEntries(entries)

	for _, e := range entries {
		if err := e.hook(ctx, xtx); err != nil {
			return err
		}
	}

	return nil
}

// RunSavepointRollbackHooks runs savepoint-scope rollback hooks for the current level. Caller must pop after.
func (r *HookRegistry) RunSavepointRollbackHooks(ctx context.Context, xtx XTransaction) {
	r.mu.RLock()

	var entries []rollbackEntry
	if len(r.savepointStack) > 0 {
		entries = make([]rollbackEntry, len(r.savepointStack[len(r.savepointStack)-1].rollback))
		copy(entries, r.savepointStack[len(r.savepointStack)-1].rollback)
	}

	r.mu.RUnlock()
	sortRollbackEntries(entries)

	for _, e := range entries {
		e.hook(ctx, xtx)
	}
}

func sortCommitEntries(entries []commitEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].order != entries[j].order {
			return entries[i].order < entries[j].order
		}

		return entries[i].index < entries[j].index
	})
}

func sortRollbackEntries(entries []rollbackEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].order != entries[j].order {
			return entries[i].order < entries[j].order
		}

		return entries[i].index < entries[j].index
	})
}

// RegisterCommit registers one or more commit hooks on the current XTransaction in ctx.
// Options apply to all hooks in this call. Defaults: ScopeTransaction, order=0.
// Returns ErrNoActiveTransaction if no XTransaction is in context.
func RegisterCommit(ctx context.Context, hooks []CommitHook, opts ...HookOption) error {
	if len(hooks) == 0 {
		return nil
	}

	xTx, ok := XTransactionFromContext(ctx)
	if !ok {
		return ErrNoActiveTransaction
	}

	reg, ok := xTx.(interface {
		RegisterCommitHooks(hooks []CommitHook, opts *HookOpts)
	})
	if !ok {
		return ErrNoActiveTransaction
	}
	o := applyOpts(opts)
	reg.RegisterCommitHooks(hooks, o)

	return nil
}

// RegisterRollback registers one or more rollback hooks on the current XTransaction in ctx.
func RegisterRollback(ctx context.Context, hooks []RollbackHook, opts ...HookOption) error {
	if len(hooks) == 0 {
		return nil
	}

	xTx, ok := XTransactionFromContext(ctx)
	if !ok {
		return ErrNoActiveTransaction
	}

	reg, ok := xTx.(interface {
		RegisterRollbackHooks(hooks []RollbackHook, opts *HookOpts)
	})
	if !ok {
		return ErrNoActiveTransaction
	}
	o := applyOpts(opts)
	reg.RegisterRollbackHooks(hooks, o)

	return nil
}
