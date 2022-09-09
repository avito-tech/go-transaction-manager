package sqlx

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-transaction-manager/internal/mock"
	"github.com/avito-tech/go-transaction-manager/transaction"
)

func Test_transactionManager_Do(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		err error
	}

	testErr := errors.New("error test")
	testCommitErr := errors.New("error Commit test")
	testRollbackErr := errors.New("error rollback test")

	tests := map[string]struct {
		args       args
		prepare    func(t *testing.T, m sqlmock.Sqlmock)
		wantErr    assert.ErrorAssertionFunc
		wantLogged []string
	}{
		"success": {
			args: args{
				ctx: context.Background(),
			},
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectCommit()
			},
			wantErr: assert.NoError,
		},
		"begin_error": {
			args: args{
				ctx: context.Background(),
			},
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(testErr)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, testErr, err)
			},
		},
		"commit_error": {
			args: args{
				ctx: context.Background(),
			},
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectCommit().WillReturnError(testCommitErr)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, testCommitErr, err)
			},
		},
		"rollback_after_error": {
			args: args{
				ctx: context.Background(),
				err: testErr,
			},
			prepare: func(t *testing.T, m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(testRollbackErr)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, testErr, err) && assert.Error(t, testRollbackErr, err)
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

			tr := NewTransactionManager(sqlx.NewDb(db, "sqlmock"), WithLog(log))

			err := tr.Do(tt.args.ctx, func(ctx context.Context) error {
				return tr.Do(ctx, func(ctx context.Context) error {
					return tt.args.err
				})
			})

			tt.wantErr(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	}
}

func Test_transactionManager_Do_Panic(t *testing.T) {
	t.Parallel()

	testPanic := "panic"
	testRollbackErr := errors.New("rollback error")

	log := mock.NewLog()
	db, dbmock, _ := sqlmock.New()

	dbmock.ExpectBegin()
	dbmock.ExpectRollback().WillReturnError(testRollbackErr)

	m := NewTransactionManager(sqlx.NewDb(db, "sqlmock"), WithLog(log))

	defer func() {
		p := recover()

		assert.Equal(t, testPanic, p)
		assert.NoError(t, dbmock.ExpectationsWereMet())
		assert.Equal(t, []string{"rollback error, panic"}, log.Logged)
	}()

	_ = m.Do(context.Background(), func(ctx context.Context) error {
		return m.Do(ctx, func(ctx context.Context) error {
			tx := TrFromCtx(ctx, nil)

			assert.NotNil(t, tx)

			panic(testPanic)
		})
	})

	assert.NoError(t, errors.New("should not be here"))
}

func Test_transactionManager_Concurrent(t *testing.T) {
	t.Parallel()

	db, dbmock, _ := sqlmock.New()

	dbmock.ExpectBegin()
	dbmock.ExpectRollback()

	var (
		wg     sync.WaitGroup
		mx     sync.Mutex
		errRet error
	)

	goroutine := func(m transaction.Manager, ctx context.Context) {
		defer wg.Done()

		err := m.Do(ctx, func(ctx context.Context) error {
			tx := TrFromCtx(ctx, nil)

			assert.NotNil(t, tx)

			select {
			case <-ctx.Done():
			case <-time.After(time.Second):
				assert.Truef(t, false, "Should not be here")
			}

			return nil
		})

		mx.Lock()
		defer mx.Unlock()

		if errRet == nil && err != nil {
			errRet = err
		}
	}

	m := NewTransactionManager(sqlx.NewDb(db, "sqlmock"), WithLog(mock.NewZeroLog()))

	err := m.Do(context.Background(), func(ctx context.Context) error {
		var err error

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go goroutine(m, ctx)
		}

		return err
	})

	wg.Wait()

	assert.NoError(t, err)
	assert.NoError(t, errRet)
}

func TestNewTransactionManager_DefaultLog(t *testing.T) {
	t.Parallel()

	type args struct {
		db *sqlx.DB
		l  logger
	}

	tests := map[string]struct {
		args args
		want func(t *testing.T, m *TrManager)
	}{
		"db_is_empty": {
			args: args{
				db: nil,
				l:  defaultLog,
			},
			want: func(t *testing.T, m *TrManager) {
				assert.Nil(t, m.db)
				assert.NotNil(t, recover())
			},
		},
		"log_is_empty": {
			args: args{
				db: &sqlx.DB{},
				l:  nil,
			},
			want: func(t *testing.T, m *TrManager) {
				assert.Equal(t, defaultLog, m.log)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var m TrManager

			defer tt.want(t, &m)

			m = *NewTransactionManager(tt.args.db, WithLog(tt.args.l)) //nolint:forcetypeassert
		})
	}
}
