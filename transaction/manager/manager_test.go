package manager

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-transaction-manager/transaction"
	mock_log "github.com/avito-tech/go-transaction-manager/transaction/manager/mock"
	"github.com/avito-tech/go-transaction-manager/transaction/mock"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

func Test_transactionManager_Do(t *testing.T) {
	t.Parallel()

	type fields struct {
		factory  transaction.TrFactory
		settings transaction.Settings
		log      logger
	}

	type args struct {
		ctx      context.Context
		settings transaction.Settings
	}

	// testErr := errors.New("error test")
	// testCommitErr := errors.New("error Commit test")
	// testRollbackErr := errors.New("error rollback test")

	tests := map[string]struct {
		args       args
		fields     func(t *testing.T, ctrl *gomock.Controller, a args) fields
		ret        error
		wantErr    assert.ErrorAssertionFunc
		wantLogged []string
	}{
		"success_commit_propagation_required": {
			args: args{
				ctx:      context.Background(),
				settings: settings.New(),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context) (transaction.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().Commit().Return(nil)

						return tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			ret:     nil,
			wantErr: assert.NoError,
		},
		"success_commit_propagation_nested": {
			args: args{
				ctx: context.Background(),
				settings: settings.New(
					settings.WithPropagation(transaction.PropagationNested),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context) (transaction.Transaction, error) {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().SavePoint(gomock.Any(), a.settings).Return(txSP, nil)

						txSP.EXPECT().Commit().Return(nil).Times(2)

						return txSP, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			ret:     nil,
			wantErr: assert.NoError,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := tt.fields(t, ctrl, tt.args)

			tr := New(f.factory, WithLog(f.log), WithSettings(f.settings))

			err := tr.DoWithSettings(
				tt.args.ctx,
				tt.args.settings,
				func(ctx context.Context) error {
					return tr.Do(ctx, func(ctx context.Context) error {
						return tt.ret
					})
				},
			)

			tt.wantErr(t, err)
		})
	}
}

func TestManager_WithOpts(t *testing.T) {
	t.Parallel()

	t.Run("set", func(t *testing.T) {
		t.Parallel()

		m := New(nil, WithLog(l{}), WithSettings(s{}))

		assert.Equal(t, l{}, m.log)
		assert.Equal(t, s{}, m.settings)
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		m := New(nil)

		assert.Equal(t, defaultLog, m.log)
		assert.Equal(t, settings.New(), m.settings)
	})
}
