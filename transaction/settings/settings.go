// Package settings implements transaction.Settings.
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
	defaultCancelable  = false
)

type ctxKey struct{}

// Opt is type to set Settings' properties.
type Opt func(s *Settings)

// Settings is an implementation of transaction.Settings.
type Settings struct {
	ctxKey       *transaction.CtxKey
	propagation  *transaction.Propagation
	isCancelable *bool
	timeout      *time.Duration
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
func (s Settings) EnrichBy(external transaction.Settings) (res transaction.Settings) { //nolint:ireturn,nolintlint
	res = s

	if s.CtxKeyOrNil() == nil {
		res = res.SetCtxKey(external.CtxKeyOrNil())
	}

	if s.PropagationOrNil() == nil {
		res = res.SetPropagation(external.PropagationOrNil())
	}

	if s.CancelableOrNil() == nil {
		res = res.SetCancelable(external.CancelableOrNil())
	}

	if s.TimeoutOrNil() == nil {
		res = res.SetTimeout(external.TimeoutOrNil())
	}

	return res
}

// CtxKey returns transaction.CtxKey for the transaction.Transaction.
func (s Settings) CtxKey() transaction.CtxKey { //nolint:ireturn,nolintlint
	if s.ctxKey == nil {
		return DefaultCtxKey
	}

	return *s.ctxKey
}

func (s Settings) CtxKeyOrNil() *transaction.CtxKey {
	return s.ctxKey
}

func (s Settings) SetCtxKey(key *transaction.CtxKey) transaction.Settings { //nolint:ireturn,nolintlint
	return s.setCtxKey(key)
}

func (s Settings) setCtxKey(key *transaction.CtxKey) Settings {
	s.ctxKey = key

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

func (s Settings) SetPropagation(p *transaction.Propagation) transaction.Settings { //nolint:ireturn,nolintlint
	return s.setPropagation(p)
}

func (s Settings) setPropagation(p *transaction.Propagation) Settings {
	s.propagation = p

	return s
}

// Cancelable defines that parent transaction.Transaction can cancel children transactions.
func (s Settings) Cancelable() bool {
	if s.isCancelable == nil {
		return defaultCancelable
	}

	return *s.isCancelable
}

func (s Settings) CancelableOrNil() *bool {
	return s.isCancelable
}

func (s Settings) SetCancelable(t *bool) transaction.Settings { //nolint:ireturn,nolintlint
	return s.setCancelable(t)
}

func (s Settings) setCancelable(t *bool) Settings {
	s.isCancelable = t

	return s
}

// TimeoutOrNil returns time.Duration of the transaction.Transaction.
func (s Settings) TimeoutOrNil() *time.Duration {
	return s.timeout
}

func (s Settings) SetTimeout(t *time.Duration) transaction.Settings { //nolint:ireturn,nolintlint
	return s.setTimeout(t)
}

func (s Settings) setTimeout(t *time.Duration) Settings {
	s.timeout = t

	return s
}
