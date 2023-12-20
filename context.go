package trm

import (
	"context"
)

// CtxKey is a type to identify trm.Transaction in a context.Context.
type CtxKey interface{}

// CtxGetter gets Transaction from context.Context.
type CtxGetter func(ctx context.Context) Transaction

// СtxManager sets and gets a Transaction in/from context.Context.
type СtxManager interface {
	// Default gets Transaction from context.Context by default CtxKey.
	Default(ctx context.Context) Transaction
	// SetDefault sets.Transaction in context.Context by default CtxKey.
	SetDefault(ctx context.Context, t Transaction) context.Context

	// ByKey gets Transaction from context.Context by CtxKey.
	ByKey(ctx context.Context, key CtxKey) Transaction
	// SetByKey sets Transaction in context.Context by.CtxKey.
	SetByKey(ctx context.Context, key CtxKey, t Transaction) context.Context
}
