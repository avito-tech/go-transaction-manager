package settings

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// DefaultCtxKey is a default key to store Transaction.
var DefaultCtxKey = ctxKey{} //nolint:gochecknoglobals

type ctxKey struct{}

// Opt is type to set Settings' properties.
type Opt func(s *Settings)

// Settings is an implementation of transaction.Settings.
type Settings struct {
	ctxKey      transaction.CtxKey
	isReadOnly  bool
	propagation transaction.Propagation
	timeout     time.Duration
}

// New creates Settings.
func New(oo ...Opt) Settings {
	s := &Settings{
		ctxKey:      ctxKey{},
		isReadOnly:  false,
		propagation: transaction.PropagationRequired,
	}

	for _, o := range oo {
		o(s)
	}

	return *s
}

// CtxKey returns transaction.CtxKey for the transaction.Transaction.
//
//nolint:ireturn,nolintlint
func (s Settings) CtxKey() transaction.CtxKey { //nolint:ireturn,nolintlint
	return s.ctxKey
}

// IsReadOnly defined that the transaction.Transaction can or cannot write data to a database.
func (s Settings) IsReadOnly() bool {
	return s.isReadOnly
}

// Propagation returns transaction.Propagation.
func (s Settings) Propagation() transaction.Propagation {
	return s.propagation
}

// Timeout returns time.Duration of the transaction.Transaction.
func (s Settings) Timeout() time.Duration {
	return s.timeout
}
