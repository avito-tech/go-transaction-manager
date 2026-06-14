package main

import (
	"context"
	"sync"

	bench "github.com/avito-tech/go-transaction-manager/trm/v2/internal/benchmark/benchutil"
)

func keyInContext() {
	ctx := context.Background()

	key := bench.CtxKey{}
	idKey := bench.IDKey(1)
	ctx = context.WithValue(ctx, key, idKey)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go nestedKeyInContext(ctx, &wg)

	wg.Wait()
}

func nestedKeyInContext(_ context.Context, wg *sync.WaitGroup) {
	wg.Done()
}
