package goredis8

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Cmdable is an experimental interface to Watch and Unwatch keys in Transaction.
type Cmdable interface {
	redis.Cmdable
	Watch
}

// Watch is experimental functional for watching updated keys.
// See redis_test.Example_watch for example.
type Watch interface {
	Watch(ctx context.Context, keys ...string) *redis.StatusCmd
	Unwatch(ctx context.Context, keys ...string) *redis.StatusCmd
}

type tx struct {
	redis.Cmdable
	tx *redis.Tx
}

func (t *tx) Watch(ctx context.Context, keys ...string) *redis.StatusCmd {
	return t.tx.Watch(ctx, keys...)
}

func (t *tx) Unwatch(ctx context.Context, keys ...string) *redis.StatusCmd {
	return t.tx.Unwatch(ctx, keys...)
}
