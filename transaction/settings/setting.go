// Package settings implements transaction.Settings.
//
//nolint:ireturn
package settings

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// DefaultCtxKey is a default key to store Transaction.
var DefaultCtxKey = ctxKey{} //nolint:gochecknoglobals

const (
	defaultIsReadOnly  = false
	defaultPropagation = transaction.PropagationRequired
	defaultTimeout     = time.Duration(0)
)

type ctxKey struct{}

// Opt is type to set Settings' properties.
type Opt func(s *Settings)

// Settings is an implementation of transaction.Settings.
type Settings struct {
	ctxKey      *transaction.CtxKey
	isReadOnly  *bool
	propagation *transaction.Propagation
	timeout     *time.Duration
}

// New creates transaction.Settings.
func New(oo ...Opt) Settings {
	s := &Settings{}

	for _, o := range oo {
		o(s)
	}

	return *s
}

//revive:disable:exported
func (s Settings) EnrichBy(external transaction.Settings) transaction.Settings {
	if s.ctxKey == nil {
		s.SetCtxKey(external.CtxKeyOrNil())
	}

	if s.isReadOnly != nil {
		s.SetIsReadOnly(external.IsReadOnlyOrNil())
	}

	if s.propagation != nil {
		s.SetPropagation(external.PropagationOrNil())
	}

	if s.timeout != nil {
		s.SetTimeout(external.TimeoutOrNil())
	}

	return s
}

// CtxKey returns transaction.CtxKey for the transaction.Transaction.
//
//nolint:ireturn,nolintlint
func (s Settings) CtxKey() transaction.CtxKey { //nolint:ireturn,nolintlint
	if s.ctxKey == nil {
		return DefaultCtxKey
	}

	return *s.ctxKey
}

func (s Settings) CtxKeyOrNil() *transaction.CtxKey {
	return s.ctxKey
}

func (s Settings) SetCtxKey(key *transaction.CtxKey) transaction.Settings {
	s.ctxKey = key

	return s
}

// IsReadOnly defined that the transaction.Transaction can or cannot write data to a database.
func (s Settings) IsReadOnly() bool {
	if s.isReadOnly == nil {
		return defaultIsReadOnly
	}

	return *s.isReadOnly
}

func (s Settings) IsReadOnlyOrNil() *bool {
	return s.isReadOnly
}

func (s Settings) SetIsReadOnly(b *bool) transaction.Settings {
	s.isReadOnly = b

	return s
}

// Propagation returns transaction.Propagation.
func (s Settings) Propagation() transaction.Propagation {
	if s.propagation == nil {
		return defaultPropagation
	}

	return *s.propagation
}

func (s Settings) PropagationOrNil() *transaction.Propagation {
	return s.propagation
}

func (s Settings) SetPropagation(p *transaction.Propagation) transaction.Settings {
	s.propagation = p

	return s
}

// Timeout returns time.Duration of the transaction.Transaction.
func (s Settings) Timeout() time.Duration {
	if s.timeout == nil {
		return defaultTimeout
	}

	return *s.timeout
}

func (s Settings) TimeoutOrNil() *time.Duration {
	return s.timeout
}

func (s Settings) SetTimeout(t *time.Duration) transaction.Settings {
	s.timeout = t

	return s
}
