//nolint:ireturn,nolintlint // return Tr for external usage.
//revive:disable:package-comments
package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/transaction"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

// TrOrDBFromCtx returns the opened Tr from the context.Context or sql.DB.
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
		tx, ok := tr.Transaction().(*sql.Tx)

		return tx, ok
	}

	return nil, false
}
