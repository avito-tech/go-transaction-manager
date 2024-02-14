package goredis8

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

const OK = "OK"

func TestTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	ctx := context.Background()
	testErr := errors.New("error test")
	testKey := "key1"
	testValue := "value"
	testExp := time.Duration(0)

	tests := map[string]struct {
		prepare  func(t *testing.T, m redismock.ClientMock)
		args     args
		ret      error
		wantErr  assert.ErrorAssertionFunc
		wantCmds int
	}{
		"commit": {
			prepare: func(t *testing.T, m redismock.ClientMock) {
				m.ExpectWatch(testKey)
				m.ExpectTxPipeline()

				m.ExpectSet(testKey, testValue, testExp).SetVal(OK)

				m.ExpectTxPipelineExec()
			},
			args: args{
				ctx: ctx,
			},
			ret:      nil,
			wantErr:  assert.NoError,
			wantCmds: 1,
		},
		"begin_error": {
			prepare: func(t *testing.T, m redismock.ClientMock) {},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "all expectations were already fulfilled, call to cmd '[watch key1]' was not expected") &&
					assert.ErrorIs(t, err, trm.ErrBegin)
			},
		},
		"commit_error": {
			prepare: func(t *testing.T, m redismock.ClientMock) {
				m.ExpectWatch(testKey)
				m.ExpectTxPipeline()

				m.ExpectSet(testKey, testValue, testExp).SetVal(OK)

				m.ExpectTxPipelineExec().RedisNil()
			},
			args: args{
				ctx: ctx,
			},
			ret: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "redis: nil") &&
					assert.ErrorIs(t, err, trm.ErrCommit)
			},
			wantCmds: 1,
		},
		"rollback": {
			prepare: func(t *testing.T, m redismock.ClientMock) {
				m.ExpectWatch(testKey)
			},
			args: args{
				ctx: ctx,
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, rmock := redismock.NewClientMock()
			log := mock.NewLog()

			tt.prepare(t, rmock)

			s := MustSettings(settings.Must(
				settings.WithPropagation(trm.PropagationNested),
			), WithWatchKeys(testKey), WithRet(&[]redis.Cmder{}))
			m := manager.Must(
				NewDefaultFactory(db),
				manager.WithLog(log),
				manager.WithSettings(s),
			)

			var tr trm.Transaction
			err := m.Do(tt.args.ctx, func(ctx context.Context) error {
				tr = trmcontext.DefaultManager.Default(ctx)

				var trNested trm.Transaction
				err := m.Do(ctx, func(ctx context.Context) error {
					trNested = trmcontext.DefaultManager.Default(ctx)

					require.NotNil(t, trNested)

					cmd := DefaultCtxGetter.DefaultTrOrDB(ctx, nil).
						Set(ctx, testKey, testValue, testExp)
					if cmd.Err() != nil {
						return cmd.Err()
					}

					return tt.ret
				})

				if trNested != nil {
					require.True(t, trNested.IsActive())
				}

				return err
			})
			if tr != nil {
				require.False(t, tr.IsActive())
			}

			if !tt.wantErr(t, err) {
				return
			}

			assert.Len(t, s.Return(), tt.wantCmds)
			assert.NoError(t, rmock.ExpectationsWereMet())
		})
	}
}

func TestTransaction_awaitDone_byContext(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	db, rmock := redismock.NewClientMock()

	f := NewDefaultFactory(db)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()

		_, tr, err := f(ctx, settings.Must())
		require.NoError(t, err)

		cancel()

		<-ctx.Done()
		require.True(t, tr.IsActive())
		<-tr.Closed()
		require.False(t, tr.IsActive())

		require.Equal(t, context.Canceled, ctx.Err())
		err = tr.Commit(ctx)
		require.ErrorIs(t, err, redis.ErrClosed)
	}()

	wg.Wait()
	assert.NoError(t, rmock.ExpectationsWereMet())
}

// TestTransaction_awaitDone_byRollback checks goroutine leak when we close transaction manually.
func TestTransaction_awaitDone_byRollback(t *testing.T) {
	t.Parallel()

	db, rmock := redismock.NewClientMock()

	f := NewDefaultFactory(db)
	ctx, _ := context.WithCancel(context.Background()) //nolint:govet

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		_, tr, err := f(ctx, settings.Must())
		require.NoError(t, err)

		require.NoError(t, tr.Rollback(ctx))
		require.False(t, tr.IsActive())
		require.NoError(t, tr.Rollback(ctx))
	}()

	wg.Wait()
	assert.NoError(t, rmock.ExpectationsWereMet())
}
