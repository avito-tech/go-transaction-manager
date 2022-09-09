//nolint:ireturn,nolintlint // return interface for external usage
//revive:disable:package-comments
package sqlx

import (
	"context"
)

type transactionKey struct{}

// ctxWithTr puts Tr in context.Context.
func ctxWithTr(ctx context.Context, tr Tr) context.Context {
	return context.WithValue(ctx, transactionKey{}, tr)
}

// TrFromCtx returns the opened Tr from the context.Context.
func TrFromCtx(ctx context.Context, db Tr) Tr {
	if tr := openedTrFromCtx(ctx); tr != nil {
		return tr
	}

	// To use sqlx.DB if the transaction was not begun.
	if db != nil {
		return db
	}

	return nil
}

// IsTrOpened checks if the transaction is open in the context.Context.
func IsTrOpened(ctx context.Context) bool {
	return openedTrFromCtx(ctx) != nil
}

func openedTrFromCtx(ctx context.Context) Tr {
	if tr, ok := ctx.Value(transactionKey{}).(Tr); ok {
		return tr
	}

	return nil
}
