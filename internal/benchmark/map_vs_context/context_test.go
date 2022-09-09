package map_vs_context

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/avito-tech/go-transaction-manager/internal/benchmark/common"
)

const (
	nestedCalls = 1
)

func BenchmarkContextEmptyTransaction(b *testing.B) {
	creator := creatorEmpty()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkContext(creator)
		}
	})
}

func BenchmarkContextCopy(b *testing.B) {
	creator := creatorCopy(getDB())

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkContext(creator)
		}
	})
}

func BenchmarkContextRealTransaction(b *testing.B) {
	creator := creatorRealTransaction(getDB())

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkContext(creator)
		}
	})
}

func benchmarkContext(creator creator) {
	ctx := context.Background()

	tr := creator()

	ctx = context.WithValue(ctx, common.CtxKey{}, tr)

	contextRunNested(ctx)
}

func contextRunNested(ctx context.Context) {
	var wgNested sync.WaitGroup

	for j := 0; j < nestedCalls; j++ {
		wgNested.Add(1)

		go contextNested(ctx, &wgNested)
	}

	wgNested.Wait()
}

func contextNested(ctx context.Context, wgNested *sync.WaitGroup) {
	defer wgNested.Done()

	t := ctx.Value(common.CtxKey{}).(*sql.Tx)

	_ = t
}
