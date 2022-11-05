// Package context implement a setter and getter to put and get trm.Transaction from context.Context.
//
//nolint:ireturn,nolintlint
package context

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
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

//revive:disable:exported
func (c *Manager) Default(ctx context.Context) trm.Transaction {
	return c.ByKey(ctx, c.ctxKey)
}

func (c *Manager) SetDefault(ctx context.Context, t trm.Transaction) context.Context {
	return c.SetByKey(ctx, c.ctxKey, t)
}

func (c *Manager) ByKey(ctx context.Context, key trm.CtxKey) trm.Transaction {
	if tr, ok := ctx.Value(key).(trm.Transaction); ok {
		return tr
	}

	return nil
}

func (c *Manager) SetByKey(ctx context.Context, key trm.CtxKey, t trm.Transaction) context.Context {
	return context.WithValue(ctx, key, t)
}
