package manager

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers/mock"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

func TestChainedMW_Do(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		fn  func(ctx context.Context) error
	}

	ctxSource := context.Background()
	ctx1LVL := context.WithValue(ctxSource, "k1", "v1")
	ctx2LVL := context.WithValue(ctxSource, "k2", "v2")

	tests := map[string]struct {
		mm      func(t *testing.T, ctrl *gomock.Controller) []trm.Manager
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"empty": {
			mm: func(t *testing.T, ctrl *gomock.Controller) []trm.Manager {
				return nil
			},
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					require.Equal(t, ctxSource, ctx)

					return nil
				},
			},
			wantErr: assert.NoError,
		},
		"one": {
			mm: func(t *testing.T, ctrl *gomock.Controller) []trm.Manager {
				m := mock.NewMockManager(ctrl)

				m.EXPECT().Do(ctxSource, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						require.Equal(t, ctxSource, ctx)

						return fn(ctx1LVL)
					})

				return []trm.Manager{m}
			},
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					require.Equal(t, ctx1LVL, ctx)

					return nil
				},
			},
			wantErr: assert.NoError,
		},
		"two": {
			mm: func(t *testing.T, ctrl *gomock.Controller) []trm.Manager {
				m1LVL := mock.NewMockManager(ctrl)
				m1LVL.EXPECT().Do(ctxSource, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx1LVL)
					})

				m2LVL := mock.NewMockManager(ctrl)
				m2LVL.EXPECT().Do(ctx1LVL, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx2LVL)
					})

				return []trm.Manager{m1LVL, m2LVL}
			},
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					require.Equal(t, ctx2LVL, ctx)

					return nil
				},
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

			c := MustChained(tt.mm(t, ctrl))

			tt.wantErr(t, c.Do(tt.args.ctx, tt.args.fn))
		})
	}
}

func TestChainedMW_DoWithSettings(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		settings trm.Settings
		fn       func(ctx context.Context) error
	}

	s := settings.Must()
	ctxSource := context.Background()
	ctx1LVL := context.WithValue(ctxSource, "k1", "v1")
	ctx2LVL := context.WithValue(ctxSource, "k2", "v2")

	tests := map[string]struct {
		mm      func(t *testing.T, ctrl *gomock.Controller) []trm.Manager
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		"empty": {
			mm: func(t *testing.T, ctrl *gomock.Controller) []trm.Manager {
				return nil
			},
			args: args{
				ctx:      context.Background(),
				settings: s,
				fn: func(ctx context.Context) error {
					require.Equal(t, ctxSource, ctx)

					return nil
				},
			},
			wantErr: assert.NoError,
		},
		"one": {
			mm: func(t *testing.T, ctrl *gomock.Controller) []trm.Manager {
				m := mock.NewMockManager(ctrl)

				m.EXPECT().DoWithSettings(ctxSource, s, gomock.Any()).
					DoAndReturn(func(ctx context.Context, _ settings.Settings, fn func(ctx context.Context) error) error {
						require.Equal(t, ctxSource, ctx)

						return fn(ctx1LVL)
					})

				return []trm.Manager{m}
			},
			args: args{
				ctx:      context.Background(),
				settings: settings.Must(),
				fn: func(ctx context.Context) error {
					require.Equal(t, ctx1LVL, ctx)

					return nil
				},
			},
			wantErr: assert.NoError,
		},
		"two": {
			mm: func(t *testing.T, ctrl *gomock.Controller) []trm.Manager {
				m1LVL := mock.NewMockManager(ctrl)
				m1LVL.EXPECT().DoWithSettings(ctxSource, s, gomock.Any()).
					DoAndReturn(func(ctx context.Context, _ settings.Settings, fn func(ctx context.Context) error) error {
						return fn(ctx1LVL)
					})

				m2LVL := mock.NewMockManager(ctrl)
				m2LVL.EXPECT().DoWithSettings(ctx1LVL, s, gomock.Any()).
					DoAndReturn(func(ctx context.Context, _ settings.Settings, fn func(ctx context.Context) error) error {
						return fn(ctx2LVL)
					})

				return []trm.Manager{m1LVL, m2LVL}
			},
			args: args{
				ctx:      context.Background(),
				settings: s,
				fn: func(ctx context.Context) error {
					require.Equal(t, ctx2LVL, ctx)

					return nil
				},
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

			c := MustChained(tt.mm(t, ctrl))

			err := c.DoWithSettings(tt.args.ctx, tt.args.settings, tt.args.fn)

			tt.wantErr(t, err)
		})
	}
}
