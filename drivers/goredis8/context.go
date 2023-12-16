package goredis8

import (
	"context"

	"github.com/go-redis/redis/v8"

	trm "github.com/avito-tech/go-transaction-manager/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/v2/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets goredis8.Pipeliner from trm.СtxManager by casting trm.Transaction to redis.UniversalClient.
type CtxGetter struct {
	ctxManager trm.СtxManager
}

// NewCtxGetter returns *CtxGetter to get Cmdable from context.Context.
func NewCtxGetter(c trm.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

// DefaultTrOrDB returns Cmdable from context.Context or DB(goredis8.Cmdable) otherwise.
func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db redis.Cmdable) redis.Cmdable {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

// TrOrDB returns Cmdable from context.Context by trm.CtxKey or DB(goredis8.Cmdable) otherwise.
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
