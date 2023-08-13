//nolint:ireturn,nolintlint // return Tr for external usage.
//revive:disable:package-comments
package pgx

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
//
//nolint:gochecknoglobals
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Tr from trm.СtxManager by casting trm.Transaction to Tr.
type CtxGetter struct {
	ctxManager trm.СtxManager
}

//revive:disable:exported
func NewCtxGetter(c trm.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db Tr) Tr {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db Tr) Tr {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr trm.Transaction) Tr {
	if tx, ok := tr.Transaction().(pgx.Tx); ok {
		return tx
	}

	return nil
}
