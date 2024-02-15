//go:build with_real_db
// +build with_real_db

package mongo_test

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	trmmongo "github.com/avito-tech/go-transaction-manager/mongo"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

// Example demonstrates the implementation of the Repository pattern by trm.Manager.
func Example() {
	client, err := mongo.NewClient(options.Client().
		ApplyURI("mongodb://127.0.0.1:27017"))
	checkErr(err)

	ctx := context.Background()

	checkErr(client.Connect(ctx))
	defer client.Disconnect(ctx)

	collection := client.Database("test").Collection("users")

	r := newRepo(collection, trmsql.DefaultCtxGetter)

	u := &user{
		ID:       1,
		Username: "username",
	}

	trManager := manager.Must(
		trmmongo.NewDefaultFactory(client),
		manager.WithCtxManager(trmcontext.DefaultManager),
	)

	err = trManager.Do(ctx, func(ctx context.Context) error {
		if err := r.Save(ctx, u); err != nil {
			return err
		}

		return trManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			return r.Save(ctx, u)
		})
	})
	checkErr(err)

	userFromDB, err := r.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB)

	// Output: &{1 new_username}
}

type repo struct {
	collection *mongo.Collection
	getter     *trmsql.CtxGetter
}

func newRepo(collection *mongo.Collection, c *trmsql.CtxGetter) *repo {
	return &repo{
		collection: collection,
		getter:     c,
	}
}

type user struct {
	ID       int64  `bson:"_id"`
	Username string `bson:"username"`
}

func (r *repo) GetByID(ctx context.Context, id int64) (*user, error) {
	var result user

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)

	return &result, err
}

func (r *repo) Save(ctx context.Context, u *user) error {
	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": u.ID},
		bson.M{"$set": u},
		options.FindOneAndUpdate().
			SetReturnDocument(options.After).
			SetUpsert(true),
	).Err(); err != nil {
		return err
	}

	return nil
}

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
