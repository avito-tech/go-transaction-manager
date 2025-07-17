package mongov2

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	trm "github.com/avito-tech/go-transaction-manager/trm/v2"
)

// Opt is a type to configure Settings.
type Opt func(*Settings) error

// WithSessionOpts sets up options.SessionOptions for the Settings.
func WithSessionOpts(opts *options.SessionOptionsBuilder) Opt {
	return func(s *Settings) error {
		*s = s.setSessionOpts(opts)

		return nil
	}
}

// WithTransactionOpts sets up options.TransactionOptions for the Settings.
func WithTransactionOpts(opts *options.TransactionOptionsBuilder) Opt {
	return func(s *Settings) error {
		*s = s.setTransactionOpts(opts)

		return nil
	}
}

// Settings contains settings for mongo.Transaction.
type Settings struct {
	trm.Settings
	sessionOpts     *options.SessionOptionsBuilder
	transactionOpts *options.TransactionOptionsBuilder
}

// NewSettings creates Settings.
func NewSettings(trms trm.Settings, oo ...Opt) (Settings, error) {
	s := &Settings{
		Settings:        trms,
		sessionOpts:     nil,
		transactionOpts: nil,
	}

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
		if s.SessionOpts() == nil {
			s = s.setSessionOpts(external.SessionOpts())
		}

		if s.TransactionOpts() == nil {
			s = s.setTransactionOpts(external.TransactionOpts())
		}
	}

	s.Settings = s.Settings.EnrichBy(in)

	return s
}

// SessionOpts returns *options.SessionOptions for the trm.Transaction.
func (s Settings) SessionOpts() *options.SessionOptionsBuilder {
	return s.sessionOpts
}

func (s Settings) setSessionOpts(opts *options.SessionOptionsBuilder) Settings {
	s.sessionOpts = opts

	return s
}

// TransactionOpts returns trm.CtxKey for the trm.Transaction.
func (s Settings) TransactionOpts() *options.TransactionOptionsBuilder {
	return s.transactionOpts
}

func (s Settings) setTransactionOpts(opts *options.TransactionOptionsBuilder) Settings {
	s.transactionOpts = opts

	return s
}
