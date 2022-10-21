package context

import (
	"sync/atomic"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

//nolint:gochecknoglobals
var defaultKeyGenerator = NewKeyGenerator()

// Generate unique transaction.CtxKey by KeyGenerator.
//
//nolint:ireturn,nolintlint
func Generate() transaction.CtxKey {
	return defaultKeyGenerator.Generate()
}

// KeyGenerator is a generator of transaction.CtxKey.
type KeyGenerator struct {
	key *int64
}

// NewKeyGenerator creates KeyGenerator.
func NewKeyGenerator() *KeyGenerator {
	initKey := int64(1)

	return &KeyGenerator{
		key: &initKey,
	}
}

// Generate unique transaction.CtxKey.
//
//nolint:ireturn,nolintlint
func (g *KeyGenerator) Generate() transaction.CtxKey {
	defer atomic.AddInt64(g.key, 1)

	return atomic.LoadInt64(g.key)
}
