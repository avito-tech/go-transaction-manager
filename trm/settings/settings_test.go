package settings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-transaction-manager/trm"
)

func TestSettings_EnrichBy(t *testing.T) {
	t.Parallel()

	type args struct {
		external trm.Settings
	}

	tests := map[string]struct {
		settings Settings
		args     args
		want     trm.Settings
	}{
		"update_default": {
			settings: Must(),
			args: args{
				external: Must(
					WithCtxKey(1),
					WithPropagation(trm.PropagationSupports),
					WithCancelable(true),
					WithTimeout(time.Second),
				),
			},
			want: Must(
				WithCtxKey(1),
				WithPropagation(trm.PropagationSupports),
				WithCancelable(true),
				WithTimeout(time.Second),
			),
		},
		"without_update": {
			settings: Must(
				WithCtxKey(1),
				WithPropagation(trm.PropagationSupports),
				WithCancelable(true),
				WithTimeout(time.Second),
			),
			args: args{
				external: Must(
					WithCtxKey(2),
					WithPropagation(trm.PropagationNever),
					WithCancelable(false),
					WithTimeout(time.Minute),
				),
			},
			want: Must(
				WithCtxKey(1),
				WithPropagation(trm.PropagationSupports),
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
		ctxKey       *trm.CtxKey
		isReadOnly   *bool
		propagation  *trm.Propagation
		isCancelable *bool
		timeout      *time.Duration
	}

	tests := map[string]struct {
		settings Settings
		want     func() want
	}{
		"get": {
			settings: Must(
				WithCtxKey(2),
				WithPropagation(trm.PropagationRequiresNew),
				WithCancelable(true),
				WithTimeout(time.Millisecond),
			),
			want: func() want {
				ctxKey := trm.CtxKey(2)
				isReadOnly := true
				propagation := trm.PropagationRequiresNew
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
			settings: Must(),
			want: func() want {
				ctxKey := trm.CtxKey(ctxKey{})
				isReadOnly := false
				propagation := trm.PropagationRequired
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
