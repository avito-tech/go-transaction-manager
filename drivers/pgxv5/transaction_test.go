package pgxv5

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	//nolint:govet
	ctx, _ := context.WithCancel(context.Background())

	testErr := errors.New("error test")
	testCommitErr := errors.New("error Commit test")
	testRollbackErr := errors.New("error rollback test")
	spPrepare := func(_ *testing.T, m pgxmock.PgxPoolIface) {
		m.ExpectBegin()
		m.ExpectCommit()
	}

	tests := map[string]struct {
		prepare func(t *testing.T, m pgxmock.PgxPoolIface)
		args    args
		ret     error
		wantErr assert.ErrorAssertionFunc
	}{
		"success": {
			prepare: func(t *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin()

				spPrepare(t, m)

				m.ExpectCommit()
			},
			args: args{
				ctx: ctx,
			},
			ret:     nil,
			wantErr: assert.NoError,
		},
		"begin_error": {
			prepare: func(_ *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin().WillReturnError(testErr)
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrBegin)
			},
		},
		"commit_error": {
			prepare: func(t *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin()

				spPrepare(t, m)

				m.ExpectCommit().WillReturnError(testCommitErr)
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, trm.ErrCommit)
			},
		},
		"rollback_after_error": {
			prepare: func(_ *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin()

				m.ExpectBegin()
				m.ExpectRollback()

				m.ExpectRollback().WillReturnError(testRollbackErr)
			},
			args: args{
				ctx: ctx,
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testRollbackErr) &&
					assert.ErrorIs(t, err, trm.ErrRollback)
			},
		},
		"begin_savepoint_error": {
			prepare: func(_ *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin()

				m.ExpectBegin().WillReturnError(testErr)
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrBegin) &&
					assert.ErrorIs(t, err, trm.ErrNestedBegin)
			},
		},
		"commit_savepoint_error": {
			prepare: func(_ *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin()

				m.ExpectBegin()
				m.ExpectCommit().WillReturnError(testCommitErr)

				m.ExpectRollback()
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, trm.ErrCommit) &&
					assert.NotNil(t, err, trm.ErrNestedCommit)
			},
		},
		"rollback_savepoint_after_error": {
			prepare: func(_ *testing.T, m pgxmock.PgxPoolIface) {
				m.ExpectBegin()

				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(testRollbackErr)

				m.ExpectRollback()
			},
			args: args{
				ctx: ctx,
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testRollbackErr) &&
					assert.ErrorIs(t, err, trm.ErrRollback) &&
					assert.NotNil(t, err, trm.ErrNestedRollback)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dbmock, err := pgxmock.NewPool()
			require.NoError(t, err)
			log := mock.NewLog()

			tt.prepare(t, dbmock)

			s := settings.Must(
				settings.WithPropagation(trm.PropagationNested),
			)
			m := manager.Must(
				NewDefaultFactory(dbmock),
				manager.WithLog(log),
				manager.WithSettings(s),
			)

			var tr trm.Transaction
			err = m.Do(tt.args.ctx, func(ctx context.Context) error {
				tr = trmcontext.DefaultManager.Default(ctx)

				var trNested trm.Transaction
				err := m.Do(ctx, func(ctx context.Context) error {
					trNested = trmcontext.DefaultManager.Default(ctx)

					require.NotNil(t, trNested)

					return tt.ret
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

			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	}
}

// captureRollbackTx is a minimal pgx.Tx that records the context passed to Rollback.
// All other methods are intentionally left unimplemented (the embedded nil interface panics
// if called, which is fine since the test only exercises Rollback via awaitDone).
type captureRollbackTx struct {
	pgx.Tx
	mu          sync.Mutex
	capturedCtx context.Context
}

func (c *captureRollbackTx) Rollback(ctx context.Context) error {
	c.mu.Lock()
	c.capturedCtx = ctx
	c.mu.Unlock()

	return nil
}

// TestTransaction_awaitDone_rollbackCtxNotCancelled verifies that awaitDone calls
// Rollback with a non-cancelled context so that pgx does not trigger the
// "slow write timer already active" panic (jackc/pgx#2332).
func TestTransaction_awaitDone_rollbackCtxNotCancelled(t *testing.T) {
	t.Parallel()

	captureTx := &captureRollbackTx{}
	tr := newDefaultTransaction(captureTx)

	ctx, cancel := context.WithCancel(context.Background())
	go tr.awaitDone(ctx)

	cancel()
	<-tr.Closed()

	captureTx.mu.Lock()
	rollbackCtx := captureTx.capturedCtx
	captureTx.mu.Unlock()

	require.NotNil(t, rollbackCtx, "awaitDone must call Rollback")
	assert.NoError(t, rollbackCtx.Err(),
		"awaitDone called Rollback with a cancelled context — this triggers pgx 'slow write timer already active' panic")
}

func TestTransaction_awaitDone_byContext(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	dbmock, err := pgxmock.NewPool()
	require.NoError(t, err)
	dbmock.ExpectBeginTx(pgx.TxOptions{})
	dbmock.ExpectRollback()

	f := NewDefaultFactory(dbmock)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()

		_, tr, err := f(ctx, settings.Must())
		require.NoError(t, err)

		cancel()

		<-tr.Closed()
		require.False(t, tr.IsActive())

		assert.NoError(t, dbmock.ExpectationsWereMet())
	}()

	wg.Wait()
}

func TestTransaction_awaitDone_byRollback(t *testing.T) {
	t.Parallel()

	dbmock, err := pgxmock.NewPool()
	require.NoError(t, err)
	dbmock.ExpectBeginTx(pgx.TxOptions{})
	dbmock.ExpectRollback()

	f := NewDefaultFactory(dbmock)
	ctx, _ := context.WithCancel(context.Background()) //nolint:govet

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		_, tr, err := f(ctx, settings.Must())
		require.NoError(t, err)

		require.NoError(t, tr.Rollback(ctx))
		require.False(t, tr.IsActive())
	}()

	wg.Wait()
}
