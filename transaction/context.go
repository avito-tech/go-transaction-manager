package transaction

import "context"

// CtxKey is a type to identify transaction.Transaction in a context.Context.
type CtxKey interface{}

// CtxWithTr returns Transaction by CtxKey.
func CtxWithTr(ctx context.Context, key CtxKey, tr Transaction) context.Context {
	return context.WithValue(ctx, key, tr)
}

// TrFromCtx returns Transaction by CtxKey.
func TrFromCtx(ctx context.Context, key CtxKey) Transaction { //nolint:ireturn
	if tr, ok := ctx.Value(key).(Transaction); ok {
		return tr
	}

	return nil
}
