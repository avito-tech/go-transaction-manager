package main

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/avito-tech/go-transaction-manager/trm/v2/internal/benchmark/common"
)

func trInContext() {
	ctx := context.Background()

	tr := &sqlx.Tx{}

	key := common.CtxKey{}
	ctx = context.WithValue(ctx, key, tr)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go nestedTrInContext(ctx, &wg)

	wg.Wait()
}

func nestedTrInContext(_ context.Context, wg *sync.WaitGroup) {
	wg.Done()
}
