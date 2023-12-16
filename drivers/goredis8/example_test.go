//go:build with_real_db
// +build with_real_db

package goredis8_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	trmredis "github.com/avito-tech/go-transaction-manager/v2/drivers/goredis8/v2"

	trm "github.com/avito-tech/go-transaction-manager/v2"
	"github.com/avito-tech/go-transaction-manager/v2/manager"
	"github.com/avito-tech/go-transaction-manager/v2/settings"
)

// Example demonstrates the implementation of the Repository pattern by trm.Manager.
func Example() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	rdb.FlushDB(ctx)

	r := newRepo(rdb, trmredis.DefaultCtxGetter)

	u1 := &user{
		Username: "username1",
	}
	u2 := &user{
		Username: "username2",
	}

	trManager := manager.Must(
		trmredis.NewDefaultFactory(rdb),
		manager.WithSettings(trmredis.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested)),
			trmredis.WithTxDecorator(trmredis.ReadOnlyFuncWithoutTxDecorator),
		)),
	)

	err := r.Save(ctx, u1)
	checkErr(err)

	var cmds []redis.Cmder
	err = trManager.DoWithSettings(
		ctx,
		trmredis.MustSettings(settings.Must(), trmredis.WithRet(&cmds)),
		func(ctx context.Context) error {
			if err := r.Save(ctx, u2); err != nil {
				return err
			}

			u1FromDB, err := r.GetByID(ctx, u1.ID)
			checkErr(err)

			fmt.Println(u1FromDB)

			return trManager.Do(ctx, func(ctx context.Context) error {
				u2.Username = "new_username2"

				return r.Save(ctx, u2)
			})
		},
	)
	checkErr(err)

	fmt.Println("cmds:", len(cmds))

	u2FromDB, err := r.GetByID(ctx, u2.ID)
	checkErr(err)

	fmt.Println(u2FromDB)
	fmt.Println(rdb.DBSize(ctx))

	// Output: &{6f2555ba-40a9-4fc8-90da-b968fd66f2e8 username1}
	// cmds: 2
	// &{0fa1769c-6e43-11ed-a1eb-0242ac120002 new_username2}
	// dbsize: 2
}

type repo struct {
	db     redis.UniversalClient
	getter *trmredis.CtxGetter
}

func newRepo(db redis.UniversalClient, c *trmredis.CtxGetter) *repo {
	return &repo{
		db:     db,
		getter: c,
	}
}

const (
	uuid1 uuid = "6f2555ba-40a9-4fc8-90da-b968fd66f2e8"
	uuid2 uuid = "0fa1769c-6e43-11ed-a1eb-0242ac120002"
)

var currentUUID = 0

func newUUID() uuid {
	res := []uuid{uuid1, uuid2}[currentUUID]

	currentUUID++

	return res
}

type uuid string

type user struct {
	ID       uuid
	Username string
}

type userRecord struct {
	Username string `redis:"username"`
}

func (r userRecord) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

func (r *userRecord) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *repo) GetByID(ctx context.Context, id uuid) (*user, error) {
	tx := r.getter.DefaultTrOrDB(ctx, r.db)
	cmd := tx.Get(ctx, string(id))

	var record userRecord

	err := cmd.Scan(&record)
	if err != nil {
		return nil, err
	}

	return r.toModel(id, record), nil
}

func (r *repo) Save(ctx context.Context, u *user) error {
	if u.ID == "" {
		u.ID = newUUID()
	}

	cmd := r.getter.DefaultTrOrDB(ctx, r.db).Set(
		ctx,
		string(u.ID),
		r.toRecord(u),
		0,
	)

	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *repo) toRecord(model *user) userRecord {
	return userRecord{
		Username: model.Username,
	}
}

func (r *repo) toModel(id uuid, row userRecord) *user {
	return &user{
		ID:       id,
		Username: row.Username,
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
