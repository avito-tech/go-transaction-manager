package mongov2

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
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
			settings: MustSettings(settings.Must()),
			args: args{
				external: MustSettings(
					settings.Must(settings.WithCancelable(true)),
					WithSessionOpts((&options.SessionOptionsBuilder{}).
						SetCausalConsistency(true)),
					WithTransactionOpts((&options.TransactionOptionsBuilder{}).
						SetReadConcern(readconcern.Majority())),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptionsBuilder{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptionsBuilder{}).
					SetReadConcern(readconcern.Majority())),
			),
		},
		"without_update": {
			settings: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptionsBuilder{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptionsBuilder{}).
					SetReadConcern(readconcern.Majority())),
			),
			args: args{
				external: MustSettings(
					settings.Must(settings.WithCancelable(false)),
					WithSessionOpts((&options.SessionOptionsBuilder{}).
						SetCausalConsistency(false)),
					WithTransactionOpts((&options.TransactionOptionsBuilder{}).
						SetReadConcern(readconcern.Local())),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptionsBuilder{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptionsBuilder{}).
					SetReadConcern(readconcern.Majority())),
			),
		},
		"update_only_trm.Settings": {
			settings: MustSettings(
				settings.Must(),
				WithSessionOpts((&options.SessionOptionsBuilder{}).
					SetCausalConsistency(true)),
			),
			args: args{
				external: settings.Must(settings.WithCancelable(true)),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptionsBuilder{}).
					SetCausalConsistency(true)),
			),
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.settings.EnrichBy(tt.args.external)

			assert.Equal(t, tt.want.CtxKey(), got.CtxKey())
			assert.Equal(t, tt.want.Propagation(), got.Propagation())
			assert.Equal(t, tt.want.Cancelable(), got.Cancelable())
			assert.Equal(t, tt.want.TimeoutOrNil(), got.TimeoutOrNil())

			assert.Equal(t, len(tt.want.(Settings).SessionOpts().List()), len(got.(Settings).SessionOpts().List()))
		})
	}
}
