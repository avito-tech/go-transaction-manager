package sql

import (
	"context"
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Tr from trm.СtxManager by casting trm.Transaction to Tr.
type CtxGetter struct {
	ctxManager trm.CtxManager
}

// NewCtxGetter returns *CtxGetter to get Tr from context.Context.
func NewCtxGetter(c trm.CtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

// DefaultTrOrDB returns Tr from context.Context or DB(Tr) otherwise.
func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db Tr) Tr {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

// TrOrDB returns Tr from context.Context by trm.CtxKey or DB(Tr) otherwise.
func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db Tr) Tr {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr trm.Transaction) Tr {
	if tx, ok := tr.Transaction().(*sql.Tx); ok {
		return tx
	}

	return nil
}
