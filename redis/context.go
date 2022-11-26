//nolint:ireturn,nolintlint // return Tr for external usage.
//revive:disable:package-comments
package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
//
//nolint:gochecknoglobals
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets redis.Pipeliner from trm.СtxManager by casting trm.Transaction to redis.UniversalClient.
type CtxGetter struct {
	ctxManager trm.СtxManager
}

//revive:disable:exported
func NewCtxGetter(c trm.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db redis.Cmdable) redis.Cmdable {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db redis.Cmdable) redis.Cmdable {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr trm.Transaction) Cmdable {
	if tx, ok := tr.Transaction().(Cmdable); ok {
		return tx
	}

	return nil
}
