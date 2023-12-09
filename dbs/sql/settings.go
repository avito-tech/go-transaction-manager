package sql

import (
	"database/sql"

trm "github.com/avito-tech/go-transaction-manager/v2"
)

// Opt is a type to configure Settings.
type Opt func(*Settings) error

// WithTxOptions sets up sql.TxOptions for the Settings.
func WithTxOptions(opts *sql.TxOptions) Opt {
	return func(s *Settings) error {
		*s = s.setTrOpts(opts)

		return nil
	}
}

// Settings contains settings for sql.Transaction.
type Settings struct {
	trm.Settings
	txOpts *sql.TxOptions
}

// NewSettings creates Settings.
func NewSettings(trms trm.Settings, oo ...Opt) (Settings, error) {
	s := &Settings{Settings: trms, txOpts: nil}

	for _, o := range oo {
		if err := o(s); err != nil {
			return Settings{}, err
		}
	}

	return *s, nil
}

// MustSettings returns Settings if err is nil and panics otherwise.
func MustSettings(trms trm.Settings, oo ...Opt) Settings {
	s, err := NewSettings(trms, oo...)
	if err != nil {
		panic(err)
	}

	return s
}

// EnrichBy fills nil properties from external Settings.
func (s Settings) EnrichBy(in trm.Settings) trm.Settings {
	external, ok := in.(Settings)
	if ok {
		if s.TxOpts() == nil {
			s = s.setTrOpts(external.TxOpts())
		}
	}

	s.Settings = s.Settings.EnrichBy(in)

	return s
}

// TxOpts returns trm.CtxKey for the trm.Transaction.
func (s Settings) TxOpts() *sql.TxOptions {
	return s.txOpts
}

func (s Settings) setTrOpts(opts *sql.TxOptions) Settings {
	s.txOpts = opts

	return s
}
