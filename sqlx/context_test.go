package sqlx

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTrFromCtx(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		db  *sqlx.DB
	}

	tests := map[string]struct {
		args args
		want Tr
	}{
		"Tx": {
			args: args{
				ctx: ctxWithTr(context.Background(), &sqlx.Tx{}),
				db:  &sqlx.DB{},
			},
			want: &sqlx.Tx{},
		},
		"DB": {
			args: args{
				ctx: context.Background(),
				db:  &sqlx.DB{},
			},
			want: &sqlx.DB{},
		},
		"nil": {
			args: args{
				ctx: context.Background(),
				db:  nil,
			},
			want: (*sqlx.DB)(nil),
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, TrFromCtx(tt.args.ctx, tt.args.db), "TrFromCtx(%v, %v)", tt.args.ctx, tt.args.db)
		})
	}
}

func TestIsTrOpened(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	tests := map[string]struct {
		args args
		want bool
	}{
		"opened": {
			args: args{
				ctx: ctxWithTr(context.Background(), &sqlx.Tx{}),
			},
			want: true,
		},
		"nil": {
			args: args{
				ctx: context.Background(),
			},
			want: false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(t, tt.want, IsTrOpened(tt.args.ctx), "IsTrOpened(%v)", tt.args.ctx)
		})
	}
}
