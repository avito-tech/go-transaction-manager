package sql

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/internal/mock"
	"github.com/avito-tech/go-transaction-manager/internal/test"
	"github.com/avito-tech/go-transaction-manager/trm"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
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
	spPrepare := func(_ *testing.T, m sqlmock.Sqlmock) {
		m.ExpectExec("SAVEPOINT tx_1").
			WillReturnResult(sqlmock.NewResult(0, 0))
		m.ExpectExec("RELEASE SAVEPOINT tx_1").
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	tests := map[string]struct {
		prepare func(t *testing.T, m sqlmock.Sqlmock)
		args    args
		ret     error
		wantErr assert.ErrorAssertionFunc
	}{
		"success": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
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
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(testErr)
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrBegin)
			},
		},
		"commit_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				spPrepare(t, m)

				m.ExpectCommit().WillReturnError(testCommitErr)
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, trm.ErrCommit)
			},
		},
		"rollback_after_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec("SAVEPOINT tx_1").
					WillReturnResult(sqlmock.NewResult(0, 0))
				m.ExpectExec("ROLLBACK TO SAVEPOINT tx_1").
					WillReturnResult(sqlmock.NewResult(0, 0))

				m.ExpectRollback().WillReturnError(testRollbackErr)
			},
			args: args{
				ctx: ctx,
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testRollbackErr) &&
					assert.ErrorIs(t, err, trm.ErrRollback)
			},
		},
		"begin_savepoint_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec("SAVEPOINT tx_1").
					WillReturnError(testErr)
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrBegin) &&
					assert.ErrorIs(t, err, trm.ErrNestedBegin)
			},
		},
		"commit_savepoint_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec("SAVEPOINT tx_1").
					WillReturnResult(sqlmock.NewResult(0, 0))
				m.ExpectExec("RELEASE SAVEPOINT tx_1").
					WillReturnError(testCommitErr)

				m.ExpectRollback()
			},
			args: args{
				ctx: ctx,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, trm.ErrNestedCommit)
			},
		},
		"rollback_savepoint_after_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec("SAVEPOINT tx_1").
					WillReturnResult(sqlmock.NewResult(0, 0))
				m.ExpectExec("ROLLBACK TO SAVEPOINT tx_1").
					WillReturnError(testRollbackErr)

				m.ExpectRollback()
			},
			args: args{
				ctx: ctx,
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testRollbackErr) &&
					assert.ErrorIs(t, err, trm.ErrRollback) &&
					assert.ErrorIs(t, err, trm.ErrNestedRollback)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, dbmock := test.NewDBMockWithClose(t)
			log := mock.NewLog()

			tt.prepare(t, dbmock)

			s := settings.Must(
				settings.WithPropagation(trm.PropagationNested),
			)
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

			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	}
}

func TestTransaction_awaitDone_byContext(t *testing.T) {
	t.Parallel()

	db, dbmock := test.NewDBMock()
	dbmock.ExpectBegin()
	dbmock.ExpectClose()
	test.Cleanup(t, func() {
		require.NoError(t, db.Close())
	})

	f := NewDefaultFactory(db)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		_, tr, err := f(ctx, settings.Must())
		require.NoError(t, err)

		cancel()

		// Need to wait for the transaction to be closed.
		// https://github.com/golang/go/blob/go1.21.6/src/database/sql/sql.go#L2174
		<-time.After(time.Millisecond)

		<-ctx.Done()
		require.False(t, tr.IsActive())
		<-tr.Closed()
		require.False(t, tr.IsActive())

		err = tr.Commit(ctx)
		require.ErrorIs(t, err, sql.ErrTxDone)
	}()

	wg.Wait()
}

// TestTransaction_awaitDone_byRollback checks goroutine leak when we close transaction manually.
func TestTransaction_awaitDone_byRollback(t *testing.T) {
	t.Parallel()

	db, dbmock := test.NewDBMockWithClose(t)
	dbmock.ExpectBegin()
	dbmock.ExpectRollback()
	dbmock.ExpectClose()
	test.Cleanup(t, func() {
		_ = db.Close()
	})

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
		require.ErrorIs(t, tr.Rollback(ctx), sql.ErrTxDone)
	}()

	wg.Wait()
}
