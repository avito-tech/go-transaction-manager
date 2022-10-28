package settings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

func TestSettings_EnrichBy(t *testing.T) {
	t.Parallel()

	type args struct {
		external transaction.Settings
	}

	tests := map[string]struct {
		settings Settings
		args     args
		want     transaction.Settings
	}{
		"update_default": {
			settings: New(),
			args: args{
				external: New(
					WithCtxKey(1),
					WithPropagation(transaction.PropagationSupports),
					WithCancelable(true),
					WithTimeout(time.Second),
				),
			},
			want: New(
				WithCtxKey(1),
				WithPropagation(transaction.PropagationSupports),
				WithCancelable(true),
				WithTimeout(time.Second),
			),
		},
		"without_update": {
			settings: New(
				WithCtxKey(1),
				WithPropagation(transaction.PropagationSupports),
				WithCancelable(true),
				WithTimeout(time.Second),
			),
			args: args{
				external: New(
					WithCtxKey(2),
					WithPropagation(transaction.PropagationNever),
					WithCancelable(false),
					WithTimeout(time.Minute),
				),
			},
			want: New(
				WithCtxKey(1),
				WithPropagation(transaction.PropagationSupports),
				WithCancelable(true),
				WithTimeout(time.Second),
			),
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.settings.EnrichBy(tt.args.external)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSettings_Getter(t *testing.T) {
	t.Parallel()

	type want struct {
		ctxKey       *transaction.CtxKey
		isReadOnly   *bool
		propagation  *transaction.Propagation
		isCancelable *bool
		timeout      *time.Duration
	}

	tests := map[string]struct {
		settings Settings
		want     func() want
	}{
		"get": {
			settings: New(
				WithCtxKey(2),
				WithPropagation(transaction.PropagationRequiresNew),
				WithCancelable(true),
				WithTimeout(time.Millisecond),
			),
			want: func() want {
				ctxKey := transaction.CtxKey(2)
				isReadOnly := true
				propagation := transaction.PropagationRequiresNew
				isCancelable := true
				timeout := time.Millisecond

				return want{
					ctxKey:       &ctxKey,
					isReadOnly:   &isReadOnly,
					propagation:  &propagation,
					isCancelable: &isCancelable,
					timeout:      &timeout,
				}
			},
		},
		"get_default": {
			settings: New(),
			want: func() want {
				ctxKey := transaction.CtxKey(ctxKey{})
				isReadOnly := false
				propagation := transaction.PropagationRequired
				isCancelable := false

				return want{
					ctxKey:       &ctxKey,
					isReadOnly:   &isReadOnly,
					propagation:  &propagation,
					isCancelable: &isCancelable,
					timeout:      nil,
				}
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			want := tt.want()

			assert.Equal(t, *want.ctxKey, tt.settings.CtxKey())
			assert.Equal(t, *want.propagation, tt.settings.Propagation())
			assert.Equal(t, *want.isCancelable, tt.settings.Cancelable())
			assert.Equal(t, want.timeout, tt.settings.TimeoutOrNil())
		})
	}
}
