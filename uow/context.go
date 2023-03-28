//nolint:ireturn,nolintlint // return Tr for external usage.
//revive:disable:package-comments
package uow

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
)

var (
	DefaultCtxKey    = ctxKey{}
	DefaultCtxGetter = NewCtxGetter(trmcontext.New(DefaultCtxKey))
)

type ctxKey struct{}

func (ctxKey) String() string {
	return "uow.ctxKey"
}

type CtxGetter struct {
	ctxManager trm.СtxManager
}

// TODO see in another language Repository with UnitOfWork
func NewCtxGetter(c trm.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db trm.UoW) trm.UoW {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) DefaultTr(ctx context.Context) trm.UoW {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return stubUoW{}
}

func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db trm.UoW) trm.UoW {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) Tr(ctx context.Context, key trm.CtxKey) trm.UoW {
	if uow := c.ctxManager.ByKey(ctx, key); uow != nil {
		return c.convert(uow)
	}

	return stubUoW{}
}

func (c *CtxGetter) convert(tr trm.Transaction) trm.UoW {
	if uow, ok := tr.Transaction().(trm.UoW); ok {
		return uow
	}

	return nil
}
