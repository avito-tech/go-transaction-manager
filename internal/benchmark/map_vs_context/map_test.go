package map_vs_context

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/avito-tech/go-transaction-manager/internal/benchmark/common"
)

// run with -test.benchtime=1s -test.benchmem

type syncMap struct {
	v  map[common.IDKey]*sql.Tx
	mu *sync.RWMutex
}

func newSyncMap() *syncMap {
	return &syncMap{
		v:  make(map[common.IDKey]*sql.Tx, 2000),
		mu: &sync.RWMutex{},
	}
}

var trMap = newSyncMap()

func BenchmarkMapEmptyTransaction(b *testing.B) {
	creator := creatorEmpty()

	i := 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkMap(i, creator)
			i++
		}
	})
}

func BenchmarkMapCopy(b *testing.B) {
	creator := creatorCopy(getDB())

	i := 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkMap(i, creator)
			i++
		}
	})
}

func BenchmarkMapRealTransaction(b *testing.B) {
	creator := creatorRealTransaction(getDB())

	i := 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkMap(i, creator)
			i++
		}
	})
}

func benchmarkMap(i int, creator creator) {
	ctx := context.Background()

	key := common.IDKey(i)

	ctx = context.WithValue(ctx, common.CtxKey{}, key)

	tr := creator()

	trMap.mu.Lock()
	trMap.v[key] = tr
	trMap.mu.Unlock()

	mapRunNested(ctx)

	trMap.mu.Lock()
	delete(trMap.v, key)
	trMap.mu.Unlock()
}

func mapRunNested(ctx context.Context) {
	var wgNested sync.WaitGroup

	for j := 0; j < nestedCalls; j++ {
		wgNested.Add(1)

		go mapNested(ctx, &wgNested)
	}

	wgNested.Wait()
}

func mapNested(ctx context.Context, wgNested *sync.WaitGroup) {
	defer wgNested.Done()

	key := ctx.Value(common.CtxKey{}).(common.IDKey)

	trMap.mu.RLock()
	t := trMap.v[key]
	trMap.mu.RUnlock()

	_ = t
}
