// Package context implement a setter and getter to put and get trm.Transaction from context.Context.
package context

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

// DefaultManager is a trm.СtxManager with settings.DefaultCtxKey.
var DefaultManager = New(settings.DefaultCtxKey) //nolint:gochecknoglobals

// Manager implements trm.СtxManager.
type Manager struct {
	ctxKey trm.CtxKey
}

// New is a factory for Manager.
func New(ctxKey trm.CtxKey) *Manager {
	return &Manager{ctxKey: ctxKey}
}

// Default returns trm.Transaction from context.Context by default key.
func (c *Manager) Default(ctx context.Context) trm.Transaction {
	return c.ByKey(ctx, c.ctxKey)
}

// SetDefault puts trm.Transaction in context.Context by default key.
func (c *Manager) SetDefault(ctx context.Context, t trm.Transaction) context.Context {
	return c.SetByKey(ctx, c.ctxKey, t)
}

// ByKey returns trm.Transaction from context.Context by key(trm.CtxKey).
func (c *Manager) ByKey(ctx context.Context, key trm.CtxKey) trm.Transaction {
	if tr, ok := ctx.Value(key).(trm.Transaction); ok {
		return tr
	}

	return nil
}

// SetByKey puts trm.Transaction in context.Context by key(trm.CtxKey).
func (c *Manager) SetByKey(ctx context.Context, key trm.CtxKey, t trm.Transaction) context.Context {
	return context.WithValue(ctx, key, t)
}
