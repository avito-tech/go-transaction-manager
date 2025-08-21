//go:build go1.16
// +build go1.16

package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Tr from trm.СtxManager by casting trm.Transaction to *gorm.DB.
type CtxGetter struct {
	ctxManager trm.CtxManager
}

// NewCtxGetter returns *CtxGetter to get *gorm.DB from context.Context.
func NewCtxGetter(c trm.CtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

// DefaultTrOrDB returns Tr(*gorm.DB) from context.Context or DB(*gorm.DB) otherwise.
func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

// TrOrDB returns Tr(*gorm.DB) from context.Context by trm.CtxKey or DB(*gorm.DB) otherwise.
func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db *gorm.DB) *gorm.DB {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr trm.Transaction) *gorm.DB {
	if tx, ok := tr.Transaction().(*gorm.DB); ok {
		return tx
	}

	return nil
}
