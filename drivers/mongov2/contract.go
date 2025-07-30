package mongov2

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

//nolint:interfacebloat
type client interface {
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	StartSession(opts ...options.Lister[options.SessionOptions]) (*mongo.Session, error)
	Database(name string, opts ...options.Lister[options.DatabaseOptions]) *mongo.Database
	ListDatabases(ctx context.Context, filter interface{}, opts ...options.Lister[options.ListDatabasesOptions]) (mongo.ListDatabasesResult, error)
	ListDatabaseNames(ctx context.Context, filter interface{}, opts ...options.Lister[options.ListDatabasesOptions]) ([]string, error)
	UseSession(ctx context.Context, fn func(context.Context) error) error
	UseSessionWithOptions(ctx context.Context, opts *options.SessionOptionsBuilder, fn func(context.Context) error) error
	Watch(ctx context.Context, pipeline interface{}, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error)
	NumberSessionsInProgress() int
}
