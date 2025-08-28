//go:build go1.19
// +build go1.19

package pgxv5

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Tr from trm.CtxManager by casting trm.Transaction to Tr.
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
	if tx, ok := tr.Transaction().(Tr); ok {
		return tx
	}

	return nil
}
