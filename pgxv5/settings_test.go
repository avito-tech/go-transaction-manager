//go:build go1.19
// +build go1.19

package pgxv5

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

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
					WithTxOptions(pgx.TxOptions{}),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(pgx.TxOptions{}),
			),
		},
		"without_update": {
			settings: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(pgx.TxOptions{IsoLevel: pgx.Serializable}),
			),
			args: args{
				external: MustSettings(
					settings.Must(settings.WithCancelable(false)),
					WithTxOptions(pgx.TxOptions{AccessMode: pgx.ReadOnly}),
				),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(pgx.TxOptions{IsoLevel: pgx.Serializable}),
			),
		},
		"update_only_trm.Settings": {
			settings: MustSettings(
				settings.Must(),
				WithTxOptions(pgx.TxOptions{IsoLevel: pgx.Serializable}),
			),
			args: args{
				external: settings.Must(settings.WithCancelable(true)),
			},
			want: MustSettings(
				settings.Must(settings.WithCancelable(true)),
				WithTxOptions(pgx.TxOptions{IsoLevel: pgx.Serializable}),
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
