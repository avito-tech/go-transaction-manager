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

// Opt is type to set settings' properties.
type Opt func(s *settings)

type settings struct {
	ctxKey      *transaction.CtxKey
	isReadOnly  *bool
	propagation *transaction.Propagation
	timeout     *time.Duration
}

// New creates settings.
func New(oo ...Opt) transaction.Settings {
	s := &settings{}

	for _, o := range oo {
		o(s)
	}

	return *s
}

func (s settings) EnrichBy(external transaction.Settings) transaction.Settings {
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
func (s settings) CtxKey() transaction.CtxKey { //nolint:ireturn,nolintlint
	if s.ctxKey == nil {
		return DefaultCtxKey
	}

	return *s.ctxKey
}

func (s settings) CtxKeyOrNil() *transaction.CtxKey {
	return s.ctxKey
}

func (s settings) SetCtxKey(key *transaction.CtxKey) transaction.Settings {
	s.ctxKey = key

	return s
}

// IsReadOnly defined that the transaction.Transaction can or cannot write data to a database.
func (s settings) IsReadOnly() bool {
	if s.isReadOnly == nil {
		return defaultIsReadOnly
	}

	return *s.isReadOnly
}

func (s settings) IsReadOnlyOrNil() *bool {
	return s.isReadOnly
}

func (s settings) SetIsReadOnly(b *bool) transaction.Settings {
	s.isReadOnly = b

	return s
}

// Propagation returns transaction.Propagation.
func (s settings) Propagation() transaction.Propagation {
	if s.propagation == nil {
		return defaultPropagation
	}

	return *s.propagation
}

func (s settings) PropagationOrNil() *transaction.Propagation {
	return s.propagation
}

func (s settings) SetPropagation(p *transaction.Propagation) transaction.Settings {
	s.propagation = p

	return s
}

// Timeout returns time.Duration of the transaction.Transaction.
func (s settings) Timeout() time.Duration {
	if s.timeout == nil {
		return defaultTimeout
	}

	return *s.timeout
}

func (s settings) TimeoutOrNil() *time.Duration {
	return s.timeout
}

func (s settings) SetTimeout(t *time.Duration) transaction.Settings {
	s.timeout = t

	return s
}
