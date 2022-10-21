// Package context implement a setter and getter to put and get transaction.Transaction from context.Context.
//
//nolint:ireturn
package context

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/transaction"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

// DefaultManager is a transaction.СtxManager with transaction.DefaultCtxKey.
var DefaultManager = New(settings.DefaultCtxKey) //nolint:gochecknoglobals

// Manager implements transaction.СtxManager.
type Manager struct {
	ctxKey transaction.CtxKey
}

// New is a factory for Manager.
func New(ctxKey transaction.CtxKey) *Manager {
	return &Manager{ctxKey: ctxKey}
}

//revive:disable:exported
func (c *Manager) Default(ctx context.Context) transaction.Transaction {
	return c.ByKey(ctx, c.ctxKey)
}

func (c *Manager) SetDefault(ctx context.Context, t transaction.Transaction) context.Context {
	return c.SetByKey(ctx, c.ctxKey, t)
}

func (c *Manager) ByKey(ctx context.Context, key transaction.CtxKey) transaction.Transaction {
	if tr, ok := ctx.Value(key).(transaction.Transaction); ok {
		return tr
	}

	return nil
}

func (c *Manager) SetByKey(ctx context.Context, key transaction.CtxKey, t transaction.Transaction) context.Context {
	return context.WithValue(ctx, key, t)
}
