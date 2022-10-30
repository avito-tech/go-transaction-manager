package sql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

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
					WithTxOptions(&sql.TxOptions{}),
				),
			},
			want: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{}),
			),
		},
		"without_update": {
			settings: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
			),
			args: args{
				external: NewSettings(
					settings.New(settings.WithCancelable(false)),
					WithTxOptions(&sql.TxOptions{ReadOnly: true}),
				),
			},
			want: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
			),
		},
		"update_only_transaction.Settings": {
			settings: NewSettings(
				settings.New(),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
			),
			args: args{
				external: settings.New(settings.WithCancelable(true)),
			},
			want: NewSettings(
				settings.New(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
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
