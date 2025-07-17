//go:build go1.21

package mongov2

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"

	"github.com/avito-tech/go-transaction-manager/drivers/mongov2/v2/internal/mtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

type user struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
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
	doNil := func(_ *mtest.T, _ context.Context) error {
		return nil
	}
	defaultFields := func(_ *mtest.T) fields {
		return fields{
			settings: MustSettings(settings.Must(
				settings.WithPropagation(trm.PropagationRequiresNew),
			), WithSessionOpts(&options.SessionOptionsBuilder{})),
		}
	}

	mt := mtest.New(
		t,
		mtest.NewOptions().ClientType(mtest.Mock),
	)

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
			fields: func(_ *mtest.T) fields {
				return fields{
					settings: MustSettings(settings.Must(
						settings.WithPropagation(trm.PropagationNested),
					), WithSessionOpts((&options.SessionOptionsBuilder{}).
						SetSnapshot(true).
						SetCausalConsistency(true))),
				}
			},
			args: args{
				ctx: context.Background(),
			},
			do: func(mt *mtest.T, _ context.Context) error {
				require.NotNil(mt, 1, "should not be here")

				return nil
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, trm.ErrBegin)
			},
		},
		"begin_transaction_error": {
			fields: func(_ *mtest.T) fields {
				return fields{
					settings: MustSettings(settings.Must(
						settings.WithPropagation(trm.PropagationNested),
					), WithTransactionOpts((&options.TransactionOptionsBuilder{}).
						SetWriteConcern(&writeconcern.WriteConcern{W: 0}))),
				}
			},
			args: args{
				ctx: context.Background(),
			},
			do: func(mt *mtest.T, _ context.Context) error {
				require.NotNil(mt, 1, "should not be here")

				return nil
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
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
					ID: bson.NewObjectID(),
				})

				return nil
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
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
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
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

			var tr trm.Transaction
			err := m.Do(tt.args.ctx, func(ctx context.Context) error {
				tr = trmcontext.DefaultManager.Default(ctx)

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

			if tr != nil {
				require.False(t, tr.IsActive())
			}

			if !tt.wantErr(t, err) {
				return
			}
		})
	}
}

func TestTransaction_awaitDone_byContext(t *testing.T) {
	t.Parallel()

	mt := mtest.New(
		t,
		mtest.NewOptions().
			ClientType(mtest.Mock).
			ShareClient(true),
	)

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
