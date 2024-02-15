//go:build go1.16
// +build go1.16

package gorm

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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

	ctx := context.Background()
	testErr := errors.New("error test")
	testCommitErr := errors.New("error Commit test")
	testRollbackErr := errors.New("error rollback test")
	spPrepare := func(_ *testing.T, m sqlmock.Sqlmock) {
		m.ExpectExec("^SAVEPOINT sp.+$").
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

				m.ExpectExec("^SAVEPOINT sp.+$").
					WillReturnResult(sqlmock.NewResult(0, 0))
				m.ExpectExec("^ROLLBACK TO SAVEPOINT sp.+$").
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

				m.ExpectExec("^SAVEPOINT sp.+$").
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
		"rollback_savepoint_after_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec("^SAVEPOINT sp.+$").
					WillReturnResult(sqlmock.NewResult(0, 0))
				m.ExpectExec("^ROLLBACK TO SAVEPOINT sp.+$").
					WillReturnError(testRollbackErr)

				m.ExpectRollback()
			},
			args: args{
				ctx: ctx,
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.NotNil(t, err, testRollbackErr) &&
					assert.NotNil(t, err, trm.ErrRollback) &&
					assert.NotNil(t, err, trm.ErrNestedRollback)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, dbmock := test.NewDBMockWithClose(t)
			dbgorm, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)

			log := mock.NewLog()

			tt.prepare(t, dbmock)

			s := settings.Must(
				settings.WithPropagation(trm.PropagationNested),
			)
			m := manager.Must(
				NewDefaultFactory(dbgorm),
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

func TestTransaction_awaitDone_byContext(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	db, dbmock := test.NewDBMock()
	dbmock.ExpectBegin()
	dbmock.ExpectRollback()
	dbmock.ExpectClose()
	test.Cleanup(t, func() {
		require.NoError(t, db.Close())
	})

	dbgorm, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	f := NewDefaultFactory(dbgorm)
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
		require.ErrorIs(t, err, sql.ErrTxDone)
	}()

	wg.Wait()
}

// TestTransaction_awaitDone_byRollback checks goroutine leak when we close transaction manually.
func TestTransaction_awaitDone_byRollback(t *testing.T) {
	t.Parallel()

	db, dbmock := test.NewDBMock()
	dbmock.ExpectBegin()
	dbmock.ExpectRollback()
	dbmock.ExpectClose()
	test.Cleanup(t, func() {
		require.NoError(t, db.Close())
	})

	dbgorm, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	f := NewDefaultFactory(dbgorm)
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
