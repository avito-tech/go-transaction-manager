package sql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/test"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

// TestTransaction_NestedSavePoints checks two savepoints inside one transaction: the second
// one gets the next identifier, and Commit and Rollback release them innermost first, so the
// identifier of the outer savepoint is still correct after the inner one is gone.
func TestTransaction_NestedSavePoints(t *testing.T) {
	t.Parallel()

	testErr := errors.New("error test")
	ok := sqlmock.NewResult(0, 0)

	tests := map[string]struct {
		prepare func(m sqlmock.Sqlmock)
		// What the innermost function returns, and whether the middle one swallows it.
		innerErr   error
		swallowErr bool
		wantErr    assert.ErrorAssertionFunc
	}{
		"two_savepoints_released_innermost_first": {
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("SAVEPOINT tx_1").WillReturnResult(ok)
				m.ExpectExec("SAVEPOINT tx_2").WillReturnResult(ok)
				m.ExpectExec("RELEASE SAVEPOINT tx_2").WillReturnResult(ok)
				m.ExpectExec("RELEASE SAVEPOINT tx_1").WillReturnResult(ok)
				m.ExpectCommit()
			},
			wantErr: assert.NoError,
		},
		"inner_error_rolls_back_both_savepoints_then_the_transaction": {
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("SAVEPOINT tx_1").WillReturnResult(ok)
				m.ExpectExec("SAVEPOINT tx_2").WillReturnResult(ok)
				m.ExpectExec("ROLLBACK TO SAVEPOINT tx_2").WillReturnResult(ok)
				m.ExpectExec("ROLLBACK TO SAVEPOINT tx_1").WillReturnResult(ok)
				m.ExpectRollback()
			},
			innerErr: testErr,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr)
			},
		},
		"inner_rollback_keeps_the_outer_savepoint_identifier": {
			prepare: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("SAVEPOINT tx_1").WillReturnResult(ok)
				m.ExpectExec("SAVEPOINT tx_2").WillReturnResult(ok)
				m.ExpectExec("ROLLBACK TO SAVEPOINT tx_2").WillReturnResult(ok)
				// The middle level handled the error, so it commits: tx_1, not tx_2 again.
				m.ExpectExec("RELEASE SAVEPOINT tx_1").WillReturnResult(ok)
				m.ExpectCommit()
			},
			innerErr:   testErr,
			swallowErr: true,
			wantErr:    assert.NoError,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, dbmock := test.NewDBMockWithClose(t)
			tt.prepare(dbmock)

			m := manager.Must(
				NewDefaultFactory(db),
				manager.WithLog(mock.NewLog()),
				manager.WithSettings(settings.Must(
					settings.WithPropagation(trm.PropagationNested),
				)),
			)

			err := m.Do(context.Background(), func(ctx context.Context) error {
				return m.Do(ctx, func(ctx context.Context) error {
					err := m.Do(ctx, func(context.Context) error {
						return tt.innerErr
					})
					if tt.swallowErr {
						require.ErrorIs(t, err, tt.innerErr)

						return nil
					}

					return err
				})
			})

			if !tt.wantErr(t, err) {
				return
			}

			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	}
}
