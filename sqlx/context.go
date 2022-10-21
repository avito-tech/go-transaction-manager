//nolint:ireturn,nolintlint // return Tr for external usage.
//revive:disable:package-comments
package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/avito-tech/go-transaction-manager/transaction"
	trmcontext "github.com/avito-tech/go-transaction-manager/transaction/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
//
//nolint:gochecknoglobals
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Tr from transaction.СtxManager by casting transaction.Transaction to Tr.
type CtxGetter struct {
	ctxManager transaction.СtxManager
}

//revive:disable:exported
func NewCtxGetter(c transaction.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db Tr) Tr {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) TrOrDB(ctx context.Context, key transaction.CtxKey, db Tr) Tr {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr transaction.Transaction) Tr {
	if tx, ok := tr.Transaction().(*sqlx.Tx); ok {
		return tx
	}

	return nil
}
