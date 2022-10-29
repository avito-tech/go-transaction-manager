package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"

	"github.com/avito-tech/go-transaction-manager/transaction"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
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
			settings: NewSettings(settings.New()),
			args: args{
				external: NewSettings(
					settings.New(settings.WithCancelable(true)),
					WithSessionOpts((&options.SessionOptions{}).
						SetCausalConsistency(true)),
					WithTransactionOpts((&options.TransactionOptions{}).
						SetReadConcern(readconcern.Majority())),
				),
			},
			want: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptions{}).
					SetReadConcern(readconcern.Majority())),
			),
		},
		"without_update": {
			settings: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptions{}).
					SetReadConcern(readconcern.Majority())),
			),
			args: args{
				external: NewSettings(
					settings.New(settings.WithCancelable(false)),
					WithSessionOpts((&options.SessionOptions{}).
						SetCausalConsistency(false)),
					WithTransactionOpts((&options.TransactionOptions{}).
						SetReadConcern(readconcern.Local())),
				),
			},
			want: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
				WithTransactionOpts((&options.TransactionOptions{}).
					SetReadConcern(readconcern.Majority())),
			),
		},
		"update_only_transaction.Settings": {
			settings: NewSettings(
				settings.New(),
				WithSessionOpts((&options.SessionOptions{}).
					SetCausalConsistency(true)),
			),
			args: args{
				external: settings.New(settings.WithCancelable(true)),
			},
			want: NewSettings(
				settings.New(settings.WithCancelable(true)),
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
