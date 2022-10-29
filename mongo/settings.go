package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// Opt is a type to configure Settings.
type Opt func(s *Settings)

// WithSessionOpts sets up options.SessionOptions for the Settings.
func WithSessionOpts(opts *options.SessionOptions) Opt {
	return func(s *Settings) {
		*s = s.setSessionOpts(opts)
	}
}

// WithTransactionOpts sets up options.TransactionOptions for the Settings.
func WithTransactionOpts(opts *options.TransactionOptions) Opt {
	return func(s *Settings) {
		*s = s.setTransactionOpts(opts)
	}
}

// Settings contains settings for mongo.Transaction.
type Settings struct {
	transaction.Settings
	sessionOpts     *options.SessionOptions
	transactionOpts *options.TransactionOptions
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

// SessionOpts returns *options.SessionOptions for the transaction.Transaction.
func (s Settings) SessionOpts() *options.SessionOptions {
	return s.sessionOpts
}

func (s Settings) setSessionOpts(opts *options.SessionOptions) Settings {
	s.sessionOpts = opts

	return s
}

// TransactionOpts returns transaction.CtxKey for the transaction.Transaction.
func (s Settings) TransactionOpts() *options.TransactionOptions {
	return s.transactionOpts
}

func (s Settings) setTransactionOpts(opts *options.TransactionOptions) Settings {
	s.transactionOpts = opts

	return s
}
