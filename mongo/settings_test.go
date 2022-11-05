package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
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
					WithSessionOpts((&options.SessionOptions{}).
						SetCausalConsistency(true)),
					WithTransactionOpts((&options.TransactionOptions{}).
						SetReadConcern(readconcern.Majority())),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptions{}).
					SetReadConcern(readconcern.Majority())),
			),
		},
		"without_update": {
			settings: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptions{}).
					SetReadConcern(readconcern.Majority())),
			),
			args: args{
				external: MustSettings(
					settings.Must(settings.WithCancelable(false)),
					WithSessionOpts((&options.SessionOptions{}).
						SetCausalConsistency(false)),
					WithTransactionOpts((&options.TransactionOptions{}).
						SetReadConcern(readconcern.Local())),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptions{}).
					SetReadConcern(readconcern.Majority())),
			),
		},
		"update_only_trm.Settings": {
			settings: MustSettings(
				settings.Must(),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
			),
			args: args{
				external: settings.Must(settings.WithCancelable(true)),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
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
