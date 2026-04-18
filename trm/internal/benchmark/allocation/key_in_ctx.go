package main

import (
	"context"
	"sync"

	benchutil "github.com/avito-tech/go-transaction-manager/trm/v2/internal/benchmark/common"
)

func keyInContext() {
	ctx := context.Background()

	key := benchutil.CtxKey{}
	idKey := benchutil.IDKey(1)
	ctx = context.WithValue(ctx, key, idKey)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go nestedKeyInContext(ctx, &wg)

	wg.Wait()
}

func nestedKeyInContext(_ context.Context, wg *sync.WaitGroup) {
	wg.Done()
}
