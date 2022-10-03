//nolint:ireturn,nolintlint // return Tr for external usage.
//revive:disable:package-comments
package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/avito-tech/go-transaction-manager/transaction"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

// TrOrDBFromCtx returns the opened Tr from the context.Context or sqlx.DB.
// TODO add ability to set ctxKey.
func TrOrDBFromCtx(ctx context.Context, db Tr) Tr {
	if tr, ok := TrFromCtx(ctx); ok {
		return tr
	}

	return db
}

// TrFromCtx returns the opened Tr from the context.Context.
func TrFromCtx(ctx context.Context) (Tr, bool) {
	if tr := transaction.TrFromCtx(ctx, settings.DefaultCtxKey); tr != nil {
		tx, ok := tr.Transaction().(*sqlx.Tx)

		return tx, ok
	}

	return nil, false
}
