//go:build with_real_db
// +build with_real_db

package goredis8_test

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	trmredis "github.com/avito-tech/go-transaction-manager/drivers/goredis8/v2"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

// Example demonstrates the watching of updated keys.
func Example_watch() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	rdb.FlushDB(ctx)

	r := newRepo(rdb, trmredis.DefaultCtxGetter)

	u := &user{
		ID:       uuid1,
		Username: "username",
	}

	trManager := manager.Must(
		trmredis.NewDefaultFactory(rdb),
		manager.WithSettings(trmredis.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested)),
			trmredis.WithTxDecorator(newWatchDecorator),
			trmredis.WithMulti(true),
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
			cmd := trmredis.DefaultCtxGetter.DefaultTrOrDB(ctx, nil).(trmredis.Watch).
				Unwatch(ctx)
			checkErr(cmd.Err())

			// Rewrite watching key1
			rdb.Set(ctx, string(u.ID), "", 0)

			return err
		},
	)
	fmt.Println(err)

	// Output: transaction: commit; redis: transaction failed
	// <nil>
}

type watchDecoratorExample struct {
	trmredis.Cmdable
}

func newWatchDecorator(tx trmredis.Cmdable, _ redis.Cmdable) trmredis.Cmdable {
	return &watchDecoratorExample{Cmdable: tx}
}

func (w *watchDecoratorExample) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := w.Watch(ctx, key)
	if cmd.Err() != nil {
		return cmd
	}

	return w.Cmdable.Set(ctx, key, value, expiration)
}
