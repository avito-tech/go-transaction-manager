package transaction

import "context"

// CtxKey is a type for context.Context keys.
type CtxKey interface{}

func ctxWithTr(ctx context.Context, key CtxKey, tr Transaction) context.Context {
	return context.WithValue(ctx, key, tr)
}

// TrFromCtx returns Transaction.
func TrFromCtx(ctx context.Context, key CtxKey) Transaction { //nolint:ireturn
	if tr, ok := ctx.Value(key).(Transaction); ok {
		return tr
	}

	return nil
}
