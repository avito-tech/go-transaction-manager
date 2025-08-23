package trm

import (
	"context"
)

// CtxKey is a type to identify trm.Transaction in a context.Context.
type CtxKey interface{}

// CtxGetter gets Transaction from context.Context.
type CtxGetter func(ctx context.Context) Transaction

// CtxManager sets and gets a Transaction in/from context.Context.
type CtxManager interface {
	// Default gets Transaction from context.Context by default CtxKey.
	Default(ctx context.Context) Transaction
	// SetDefault sets.Transaction in context.Context by default CtxKey.
	SetDefault(ctx context.Context, t Transaction) context.Context

	// ByKey gets Transaction from context.Context by CtxKey.
	ByKey(ctx context.Context, key CtxKey) Transaction
	// SetByKey sets Transaction in context.Context by.CtxKey.
	SetByKey(ctx context.Context, key CtxKey, t Transaction) context.Context
}

// СtxManager is old name with first non-ASCII character.
// Deprecated: Type name contains first non-ASCII character.
// Type is safed in term of backward compatibility, use above CtxManager instead.
type СtxManager = CtxManager
