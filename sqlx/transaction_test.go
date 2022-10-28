package sqlx

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/internal/mock"
	"github.com/avito-tech/go-transaction-manager/transaction"
	trmcontext "github.com/avito-tech/go-transaction-manager/transaction/context"
	"github.com/avito-tech/go-transaction-manager/transaction/manager"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

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
				ctx: context.Background(),
			},
			ret:     nil,
			wantErr: assert.NoError,
		},
		"begin_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(testErr)
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, transaction.ErrBegin)
			},
		},
		"commit_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				spPrepare(t, m)

				m.ExpectCommit().WillReturnError(testCommitErr)
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, transaction.ErrCommit)
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
				ctx: context.Background(),
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testRollbackErr) &&
					assert.ErrorIs(t, err, transaction.ErrRollback)
			},
		},
		"begin_savepoint_error": {
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()

				m.ExpectExec("SAVEPOINT tx_1").
					WillReturnError(testErr)
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, transaction.ErrBegin) &&
					assert.ErrorIs(t, err, transaction.ErrSPBegin)
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
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, transaction.ErrCommit) &&
					assert.ErrorIs(t, err, transaction.ErrSPCommit)
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
				ctx: context.Background(),
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testRollbackErr) &&
					assert.ErrorIs(t, err, transaction.ErrRollback) &&
					assert.ErrorIs(t, err, transaction.ErrSPRollback)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, dbmock, _ := sqlmock.New()
			log := mock.NewLog()

			tt.prepare(t, dbmock)

			s := settings.New(
				settings.WithPropagation(transaction.PropagationNested),
			)
			m := manager.New(
				NewDefaultFactory(sqlx.NewDb(db, "sqlmock")),
				manager.WithLog(log),
				manager.WithSettings(s),
			)

			var tr Transaction
			err := m.Do(tt.args.ctx, func(ctx context.Context) error {
				var trNested transaction.Transaction
				err := m.Do(ctx, func(ctx context.Context) error {
					trNested = trmcontext.DefaultManager.Default(ctx)

					require.NotNil(t, trNested)

					return tt.ret
				})

				if trNested != nil {
					require.Equal(t, true, trNested.IsActive())
				}

				return err
			})
			require.Equal(t, false, tr.IsActive())

			if !tt.wantErr(t, err) {
				return
			}
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	}
}
