package sql

import (
	"database/sql"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// Opt is a type to configure Settings.
type Opt func(s *Settings)

// WithTxOptions sets up sql.TxOptions for the Settings.
func WithTxOptions(opts *sql.TxOptions) Opt {
	return func(s *Settings) {
		*s = s.setTrOpts(opts)
	}
}

// Settings contains settings for mongo.Transaction.
type Settings struct {
	transaction.Settings
	txOpts *sql.TxOptions
}

// NewSettings creates Settings.
func NewSettings(trms transaction.Settings, oo ...Opt) Settings {
	s := &Settings{Settings: trms}

	for _, o := range oo {
		o(s)
	}

	return *s
}

//revive:disable:exported
func (s Settings) EnrichBy(in transaction.Settings) (res transaction.Settings) { //nolint:ireturn,nolintlint
	external, ok := in.(Settings)
	if ok {
		if s.TxOpts() == nil {
			s = s.setTrOpts(external.TxOpts())
		}
	}

	s.Settings = s.Settings.EnrichBy(in)

	return s
}

// TxOpts returns transaction.CtxKey for the transaction.Transaction.
func (s Settings) TxOpts() *sql.TxOptions {
	return s.txOpts
}

func (s Settings) setTrOpts(opts *sql.TxOptions) Settings {
	s.txOpts = opts

	return s
}
