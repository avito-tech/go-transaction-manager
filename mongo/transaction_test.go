package mongo

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/goleak"

	"github.com/avito-tech/go-transaction-manager/internal/mock"
	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
)

type user struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	type fields struct {
		settings trm.Settings
	}

	testErr := errors.New("error test")
	doNil := func(mt *mtest.T, ctx context.Context) error {
		return nil
	}
	defaultFields := func(mt *mtest.T) fields {
		return fields{
			settings: MustSettings(settings.Must(
				settings.WithPropagation(trm.PropagationRequiresNew),
			), WithSessionOpts(&options.SessionOptions{})),
		}
	}

	mt := mtest.New(
		t,
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()

	tests := map[string]struct {
		fields  func(mt *mtest.T) fields
		args    args
		do      func(mt *mtest.T, ctx context.Context) error
		wantErr assert.ErrorAssertionFunc
	}{
		"success": {
			fields: defaultFields,
			args: args{
				ctx: context.Background(),
			},
			do:      doNil,
			wantErr: assert.NoError,
		},
		"begin_session_error": {
			fields: func(mt *mtest.T) fields {
				return fields{
					settings: MustSettings(settings.Must(
						settings.WithPropagation(trm.PropagationNested),
					), WithSessionOpts((&options.SessionOptions{}).
						SetSnapshot(true).
						SetCausalConsistency(true))),
				}
			},
			args: args{
				ctx: context.Background(),
			},
			do: func(mt *mtest.T, ctx context.Context) error {
				require.NotNil(mt, 1, "should not be here")

				return nil
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, trm.ErrBegin)
			},
		},
		"begin_transaction_error": {
			fields: func(mt *mtest.T) fields {
				return fields{
					settings: MustSettings(settings.Must(
						settings.WithPropagation(trm.PropagationNested),
					), WithTransactionOpts((&options.TransactionOptions{}).
						SetWriteConcern(&writeconcern.WriteConcern{W: 0}))),
				}
			},
			args: args{
				ctx: context.Background(),
			},
			do: func(mt *mtest.T, ctx context.Context) error {
				require.NotNil(mt, 1, "should not be here")

				return nil
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, trm.ErrBegin)
			},
		},
		"commit_error": {
			fields: defaultFields,
			args: args{
				ctx: context.Background(),
			},
			do: func(mt *mtest.T, ctx context.Context) error {
				_, _ = mt.Coll.InsertOne(ctx, user{
					ID: primitive.NewObjectID(),
				})

				return nil
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var divErr mongo.CommandError

				return assert.ErrorAs(t, err, &divErr) &&
					assert.ErrorIs(t, err, trm.ErrCommit)
			},
		},
		"rollback_after_error": {
			fields: defaultFields,
			args: args{
				ctx: context.Background(),
			},
			do: func(mt *mtest.T, ctx context.Context) error {
				s := mongo.SessionFromContext(ctx)

				require.NoError(mt, s.AbortTransaction(ctx))

				return testErr
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrRollback)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		mt.Run(name, func(mt *mtest.T) {
			mt.Parallel()

			log := mock.NewLog()

			f := tt.fields(mt)

			m := manager.Must(
				NewDefaultFactory(mt.Client),
				manager.WithLog(log),
				manager.WithSettings(f.settings),
			)

			var tr Transaction
			err := m.Do(tt.args.ctx, func(ctx context.Context) error {
				var trNested trm.Transaction
				err := m.Do(ctx, func(ctx context.Context) error {
					trNested = trmcontext.DefaultManager.Default(ctx)

					require.NotNil(t, trNested)

					return tt.do(mt, ctx)
				})

				if trNested != nil {
					require.False(t, trNested.IsActive())
				}

				return err
			})
			require.False(t, false, tr.IsActive())

			if !tt.wantErr(t, err) {
				return
			}
		})
	}
}

func TestTransaction_awaitDone(t *testing.T) {
	t.Parallel()

	mt := mtest.New(
		t,
		mtest.NewOptions().
			ClientType(mtest.Mock).
			ShareClient(true),
	)
	defer mt.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	f := NewDefaultFactory(mt.Client)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()

		_, tr, err := f(ctx, settings.Must())

		cancel()
		<-time.After(time.Second)

		<-ctx.Done()

		require.NoError(mt, err)
		require.False(mt, tr.IsActive())
	}()

	wg.Wait()
}
