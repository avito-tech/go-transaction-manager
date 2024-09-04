package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Tr from trm.СtxManager by casting trm.Transaction to Tr.
type CtxGetter struct {
	ctxManager trm.СtxManager
}

// NewCtxGetter returns *CtxGetter to get Tr from context.Context.
func NewCtxGetter(c trm.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

// DefaultTrOrDB returns mongo.Session from context.Context or DB(mongo.Session) otherwise.
func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db mongo.Session) mongo.Session {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

// TrOrDB returns mongo.Session from context.Context by trm.CtxKey or DB(mongo.Session) otherwise.
func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db mongo.Session) mongo.Session {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr trm.Transaction) mongo.Session {
	if tx, ok := tr.Transaction().(mongo.Session); ok {
		return tx
	}

	return nil
}
