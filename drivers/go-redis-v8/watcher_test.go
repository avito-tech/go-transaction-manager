//go:build with_real_db
// +build with_real_db

package go_redis_v8_test

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-tech/go-transaction-manager/v2/manager"
	"github.com/avito-tech/go-transaction-manager/v2/settings"
)

// Example demonstrates the watching of updated keys.
func Example_watch() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	rdb.FlushDB(ctx)

	r := newRepo(rdb, DefaultCtxGetter)

	u := &user{
		ID:       uuid1,
		Username: "username",
	}

	trManager := manager.Must(
		NewDefaultFactory(rdb),
		manager.WithSettings(MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested)),
			WithTxDecorator(newWatchDecorator),
			WithMulti(true),
		)),
	)

	err := r.Save(ctx, u)
	checkErr(err)

	err = trManager.Do(
		ctx,
		func(ctx context.Context) error {
			u.Username = "new_username"
			err = r.Save(ctx, u)

			// Rewrite watching key1
			rdb.Set(ctx, string(u.ID), "", 0)

			return err
		},
	)
	fmt.Println(err)

	err = trManager.Do(
		ctx,
		func(ctx context.Context) error {
			u.Username = "new_username"
			err = r.Save(ctx, u)

			// Unwatch keys
			cmd := DefaultCtxGetter.DefaultTrOrDB(ctx, nil).(Watch).
				Unwatch(ctx)
			checkErr(cmd.Err())

			// Rewrite watching key1
			rdb.Set(ctx, string(u.ID), "", 0)

			return err
		},
	)
	fmt.Println(err)

	// Output: transaction: commit; go-redis-v9: transaction failed
	// <nil>
}

type watchDecoratorExample struct {
	Cmdable
}

func newWatchDecorator(tx Cmdable, _ redis.Cmdable) Cmdable {
	return &watchDecoratorExample{Cmdable: tx}
}

func (w *watchDecoratorExample) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := w.Watch(ctx, key)
	if cmd.Err() != nil {
		return cmd
	}

	return w.Cmdable.Set(ctx, key, value, expiration)
}
