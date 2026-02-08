package trm

import "context"

type xCtxKey struct{}

// XTransactionFromContext returns the current XTransaction from ctx if present.
func XTransactionFromContext(ctx context.Context) (XTransaction, bool) {
	if tx, ok := ctx.Value(xCtxKey{}).(XTransaction); ok {
		return tx, true
	}

	return nil, false
}

// WithXTransaction stores xtx in ctx and returns the new context. Used by XManager.
func WithXTransaction(ctx context.Context, xtx XTransaction) context.Context {
	return context.WithValue(ctx, xCtxKey{}, xtx)
}
