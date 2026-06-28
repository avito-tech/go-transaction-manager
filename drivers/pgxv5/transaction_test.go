package pgxv5

import (
	"context"
	"errors"
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

	//nolint:govet,gosec
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
					assert.NotErrorIs(t, err, trm.ErrNestedCommit) // driver uses own nested transactions
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
					assert.NotErrorIs(t, err, trm.ErrNestedRollback) // driver uses own nested transactions
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

// TestTransaction_Rollback_withCancelledCtx verifies that Rollback with a cancelled
// context still issues the rollback and marks the transaction closed, propagating the
// context.Canceled error rather than swallowing it (jackc/pgx#2332).
func TestTransaction_Rollback_withCancelledCtx(t *testing.T) {
	t.Parallel()

	dbmock, err := pgxmock.NewPool()
	require.NoError(t, err)
	dbmock.ExpectBeginTx(pgx.TxOptions{})
	dbmock.ExpectRollback()

	f := NewDefaultFactory(dbmock)
	ctx, cancel := context.WithCancel(context.Background())

	_, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	cancel()

	require.ErrorIs(t, tr.Rollback(ctx), context.Canceled)
	require.False(t, tr.IsActive())

	assert.NoError(t, dbmock.ExpectationsWereMet())
}

func TestTransaction_Rollback(t *testing.T) {
	t.Parallel()

	dbmock, err := pgxmock.NewPool()
	require.NoError(t, err)
	dbmock.ExpectBeginTx(pgx.TxOptions{})
	dbmock.ExpectRollback()

	f := NewDefaultFactory(dbmock)

	_, tr, err := f(context.Background(), settings.Must())
	require.NoError(t, err)

	require.NoError(t, tr.Rollback(context.Background()))
	require.False(t, tr.IsActive())

	assert.NoError(t, dbmock.ExpectationsWereMet())
}
