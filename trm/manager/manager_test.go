package manager

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	mock_log "github.com/avito-tech/go-transaction-manager/trm/v2/manager/mock"
	"github.com/avito-tech/go-transaction-manager/trm/v2/mock"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

func Test_transactionManager_Do(t *testing.T) {
	t.Parallel()

	type fields struct {
		factory  trm.TrFactory
		settings trm.Settings
		log      logger
	}

	type args struct {
		ctx            context.Context
		settings       trm.Settings
		nestedSettings trm.Settings
	}

	ctxManager := trmcontext.DefaultManager

	emptyFactory := func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
		return ctx, nil, nil
	}

	tests := map[string]struct {
		args         args
		fields       func(t *testing.T, ctrl *gomock.Controller, a args) fields
		wantErr      assert.ErrorAssertionFunc
		wantEmptyCtx bool
	}{
		"PropagationRequired_success_commit": {
			args: args{
				ctx:            context.Background(),
				settings:       settings.Must(),
				nestedSettings: settings.Must(),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().IsActive().Return(true)
						tx.EXPECT().Commit(gomock.Any()).Return(nil)

						return ctx, tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: assert.NoError,
		},
		"PropagationNested_success_commit": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationNested),
				),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationNested),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().IsActive().Return(true).Times(2)
						txSP.EXPECT().Begin(gomock.Any(), a.settings).Return(ctx, txSP, nil)
						txSP.EXPECT().Commit(gomock.Any()).Return(nil).Times(2)

						return ctx, txSP, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: assert.NoError,
		},
		"PropagationsMandatory_success_commit": {
			args: args{
				ctx: ctxManager.SetByKey(
					context.Background(),
					settings.DefaultCtxKey,
					mock.NewMockTransaction(nil),
				),
				settings: settings.Must(),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationsMandatory),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().IsActive().Return(true)
						txSP.EXPECT().Begin(gomock.Any(), a.settings).Return(txSP, nil)
						txSP.EXPECT().Commit(gomock.Any()).Return(nil).Times(2)

						return ctx, txSP, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: assert.NoError,
		},
		"PropagationsMandatory_error_without_open_transaction": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationsMandatory),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory:  emptyFactory,
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, trm.ErrPropagationMandatory)
			},
		},
		"PropagationNever_success": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(settings.WithPropagation(
					trm.PropagationNever,
				)),
				nestedSettings: settings.Must(settings.WithPropagation(
					trm.PropagationNever,
				)),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory:  emptyFactory,
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantEmptyCtx: true,
			wantErr:      assert.NoError,
		},
		"PropagationNever_error_transaction_is_opened": {
			args: args{
				ctx:      context.Background(),
				settings: settings.Must(),
				nestedSettings: settings.Must(settings.WithPropagation(
					trm.PropagationNever,
				)),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().IsActive().Return(true)
						tx.EXPECT().Rollback(gomock.Any()).Return(nil)

						return ctx, tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, trm.ErrPropagationNever)
			},
		},
		"PropagationNotSupported_success_commit": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationNotSupported),
				),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationNotSupported),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory:  emptyFactory,
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr:      assert.NoError,
			wantEmptyCtx: true,
		},
		"PropagationNotSupported_nilling_ctx_success_commit": {
			args: args{
				ctx:      context.Background(),
				settings: settings.Must(),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationNotSupported),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().IsActive().Return(true)
						txSP.EXPECT().Commit(gomock.Any()).Return(nil)

						return ctx, txSP, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr:      assert.NoError,
			wantEmptyCtx: true,
		},
		"PropagationRequiresNew": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationRequiresNew),
				),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationRequiresNew),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func() trm.TrFactory {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().IsActive().Return(true).Times(2)
						txSP.EXPECT().Commit(gomock.Any()).Return(nil).Times(2)

						return func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
							return ctx, txSP, nil
						}
					}(),
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: assert.NoError,
		},
		"PropagationSupports_nil_transaction": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationSupports),
				),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationSupports),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory:  emptyFactory,
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr:      assert.NoError,
			wantEmptyCtx: true,
		},
		"PropagationSupports": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationRequired),
				),
				nestedSettings: settings.Must(
					settings.WithPropagation(trm.PropagationSupports),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func() func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().IsActive().Return(true)
						txSP.EXPECT().Commit(gomock.Any()).Return(nil)

						return func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
							return ctx, txSP, nil
						}
					}(),
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
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

			tr := Must(f.factory, WithLog(f.log), WithSettings(f.settings))

			err := tr.DoWithSettings(
				tt.args.ctx,
				tt.args.settings,
				func(ctx context.Context) error {
					return tr.DoWithSettings(
						ctx,
						tt.args.nestedSettings,
						func(ctx context.Context) error {
							if tt.wantEmptyCtx {
								require.Nil(t, ctxManager.Default(ctx))
							} else {
								require.NotNil(t, ctxManager.Default(ctx))
							}

							return nil
						},
					)
				},
			)

			tt.wantErr(t, err)
		})
	}
}

func Test_transactionManager_Do_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		factory  trm.TrFactory
		settings trm.Settings
		log      logger
	}

	type args struct {
		ctx      context.Context
		settings trm.Settings
	}

	testErr := errors.New("error test")
	testCommitErr := errors.New("error Commit test")
	testRollbackErr := errors.New("error rollback test")
	defaultArgs := args{
		ctx:      context.Background(),
		settings: settings.Must(),
	}

	tests := map[string]struct {
		args    args
		fields  func(t *testing.T, ctrl *gomock.Controller, a args) fields
		ret     error
		wantErr assert.ErrorAssertionFunc
	}{
		//nolint:dupl
		"transaction_factory_&_rollback_error": {
			args: defaultArgs,
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				return fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().
							IsActive().
							Return(true)
						tx.EXPECT().
							Rollback(gomock.Any()).
							Return(testRollbackErr)

						return ctx, tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}
			},
			ret: testErr,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrRollback)
			},
		},
		"skip_rollback_with_error": {
			args: defaultArgs,
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				return fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().
							IsActive().
							Return(true)
						tx.EXPECT().
							Commit(gomock.Any())

						return ctx, tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}
			},
			ret: trm.Skippable(testErr),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.True(t, trm.IsSkippable(err))
			},
		},
		"skip_rollback_with_commit_error": {
			args: defaultArgs,
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				return fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().
							IsActive().
							Return(true)
						tx.EXPECT().
							Commit(gomock.Any()).
							Return(testCommitErr)

						return ctx, tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}
			},
			ret: trm.Skippable(testErr),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, trm.ErrCommit) &&
					assert.False(t, trm.IsSkippable(err))
			},
		},
		//nolint:dupl
		"commit_error": {
			args: defaultArgs,
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				return fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						tx := mock.NewMockTransaction(ctrl)

						tx.EXPECT().
							IsActive().
							Return(true)
						tx.EXPECT().
							Commit(gomock.Any()).
							Return(testCommitErr)

						return ctx, tx, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}
			},
			ret: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testCommitErr) &&
					assert.ErrorIs(t, err, trm.ErrCommit)
			},
		},
		"PropagationNested_err_nested_begin": {
			args: args{
				ctx: context.Background(),
				settings: settings.Must(
					settings.WithPropagation(trm.PropagationNested),
				),
			},
			fields: func(t *testing.T, ctrl *gomock.Controller, a args) fields {
				f := fields{
					factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
						txSP := mock.NewMocktransactionWithSP(ctrl)

						txSP.EXPECT().IsActive().Return(true).Times(1)
						txSP.EXPECT().Begin(gomock.Any(), a.settings).Return(ctx, nil, testErr)
						txSP.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)

						return ctx, txSP, nil
					},
					settings: a.settings,
					log:      mock_log.NewMocklogger(ctrl),
				}

				return f
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testErr) &&
					assert.ErrorIs(t, err, trm.ErrNestedBegin)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := tt.fields(t, ctrl, tt.args)

			tr := Must(f.factory, WithLog(f.log), WithSettings(f.settings))

			err := tr.DoWithSettings(
				tt.args.ctx,
				tt.args.settings,
				func(ctx context.Context) error {
					return tr.Do(
						ctx,
						func(ctx context.Context) error {
							return tt.ret
						},
					)
				},
			)

			tt.wantErr(t, err)
		})
	}
}

func Test_transactionManager_Do_Panic(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPanic := "panic"
	testRollbackErr := errors.New("rollback error")

	log := mock.NewLog()
	factory := func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
		tx := mock.NewMockTransaction(ctrl)

		tx.EXPECT().IsActive().Return(true)
		tx.EXPECT().Rollback(gomock.Any()).Return(testRollbackErr)

		return ctx, tx, nil
	}

	m := Must(factory, WithLog(log))
	ctxManager := trmcontext.DefaultManager

	defer func() {
		p := recover()

		assert.Equal(t, testPanic, p)
		assert.Equal(t, []string{"rollback error, panic"}, log.Logged)
	}()

	_ = m.Do(context.Background(), func(ctx context.Context) error {
		return m.Do(ctx, func(ctx context.Context) error {
			assert.NotNil(
				t,
				ctxManager.Default(ctx),
			)

			panic(testPanic)
		})
	})

	assert.NoError(t, errors.New("should not be here"))
}

//nolint:tparallel // there is not t.Cleanup in go 1.13 and less.
func Test_transactionManager_Do_ClosedTransaction(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testErr := errors.New("test error")

	tests := map[string]struct {
		ret     error
		wantErr require.ErrorAssertionFunc
	}{
		"without_error": {
			ret: nil,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.Equal(t, trm.ErrAlreadyClosed, err)
			},
		},
		"with_error": {
			ret: testErr,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, testErr)
				require.ErrorIs(t, err, trm.ErrAlreadyClosed)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tx := mock.NewMockTransaction(ctrl)
			tx.EXPECT().IsActive().Return(false).MinTimes(2)

			factory := func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
				return ctx, tx, nil
			}

			m := Must(
				factory,
				WithSettings(settings.Must(settings.WithPropagation(trm.PropagationRequiresNew))),
			)

			err := m.Do(context.Background(), func(ctx context.Context) error {
				return m.Do(ctx, func(ctx context.Context) error {
					return tt.ret
				})
			})

			tt.wantErr(t, err)
		})
	}
}

//nolint:tparallel // there is not t.Cleanup in go 1.13 and less.
func Test_transactionManager_Do_Cancel(t *testing.T) {
	type fields struct {
		settings trm.Settings
		factory  trm.TrFactory
	}

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := map[string]struct {
		fields  fields
		ctx     func(ctx context.Context) (context.Context, context.CancelFunc)
		do      func(t *testing.T, ctx context.Context)
		wantErr require.ErrorAssertionFunc
	}{
		"cancel": {
			fields: fields{
				factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
					tr := mock.NewMockTransaction(ctrl)
					tr.EXPECT().IsActive().Return(false)

					return ctx, tr, nil
				},
				settings: settings.Must(
					settings.WithCancelable(true),
					settings.WithPropagation(trm.PropagationRequiresNew),
				),
			},
			ctx: context.WithCancel,
			do: func(t *testing.T, ctx context.Context) {
				time.Sleep(time.Millisecond)

				assert.ErrorIs(t, ctx.Err(), context.Canceled)
			},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
		"timeout": {
			fields: fields{
				factory: func(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) {
					tr := mock.NewMockTransaction(ctrl)
					tr.EXPECT().IsActive().Return(false)

					return ctx, tr, nil
				},
				settings: settings.Must(
					settings.WithCancelable(true),
					settings.WithTimeout(time.Millisecond),
					settings.WithPropagation(trm.PropagationRequiresNew),
				),
			},
			ctx: func(ctx context.Context) (context.Context, context.CancelFunc) {
				return ctx, func() {}
			},
			do: func(t *testing.T, ctx context.Context) {
				select {
				case <-time.After(time.Second):
				case <-ctx.Done():
				}

				assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
			},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := Must(
				tt.fields.factory,
				WithSettings(tt.fields.settings),
			)

			wg := sync.WaitGroup{}
			var err error

			wg.Add(1)

			ctx, cancel := tt.ctx(context.Background())
			go func() {
				err = m.Do(ctx, func(ctx context.Context) error {
					return m.Do(ctx, func(ctx context.Context) error {
						tt.do(t, ctx)

						return nil
					})
				})

				wg.Done()
			}()

			cancel()

			wg.Wait()

			tt.wantErr(t, err)
		})
	}
}

func TestManager_WithOpts(t *testing.T) {
	t.Parallel()

	t.Run("set", func(t *testing.T) {
		t.Parallel()

		l := mock.NewZeroLog()
		m := Must(nil, WithLog(l), WithSettings(s{}))

		assert.Equal(t, l, m.log)
		assert.Equal(t, s{}, m.settings)
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		m := Must(nil)

		assert.Equal(t, defaultLog, m.log)
		assert.Equal(t, settings.Must(), m.settings)
	})
}
