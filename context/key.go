package context

import (
	"sync/atomic"

	trm "github.com/avito-tech/go-transaction-manager/v2"
)

//nolint:gochecknoglobals
var defaultKeyGenerator = NewKeyGenerator()

// Generate unique trm.CtxKey by KeyGenerator.
//
//nolint:ireturn,nolintlint
func Generate() trm.CtxKey {
	return defaultKeyGenerator.Generate()
}

// KeyGenerator is a generator of trm.CtxKey.
type KeyGenerator struct {
	key int64
}

// NewKeyGenerator creates KeyGenerator.
func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{
		key: 1,
	}
}

// Generate unique trm.CtxKey.
//
//nolint:ireturn,nolintlint
func (g *KeyGenerator) Generate() trm.CtxKey {
	defer atomic.AddInt64(&g.key, 1)

	return atomic.LoadInt64(&g.key)
}
