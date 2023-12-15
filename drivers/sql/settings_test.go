package sql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	trm "github.com/avito-tech/go-transaction-manager/v2"
	"github.com/avito-tech/go-transaction-manager/v2/settings"
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
					WithTxOptions(&sql.TxOptions{}),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{}),
			),
		},
		"without_update": {
			settings: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
			),
			args: args{
				external: MustSettings(
					settings.Must(settings.WithCancelable(false)),
					WithTxOptions(&sql.TxOptions{ReadOnly: true}),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
			),
		},
		"update_only_trm.Settings": {
			settings: MustSettings(
				settings.Must(),
				WithTxOptions(&sql.TxOptions{Isolation: sql.LevelWriteCommitted}),
			),
			args: args{
				external: settings.Must(settings.WithCancelable(true)),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
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
