// Package settings implements trm.Settings.
package settings

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

// DefaultCtxKey is a default key to store Transaction.
var DefaultCtxKey = ctxKey{} //nolint:gochecknoglobals

const (
	defaultPropagation = trm.PropagationRequired
	defaultCancelable  = false
)

type ctxKey struct{}

// Opt is type to set Settings' properties.
type Opt func(*Settings) error

// Settings is an implementation of trm.Settings.
type Settings struct {
	ctxKey       *trm.CtxKey
	propagation  *trm.Propagation
	isCancelable *bool
	timeout      *time.Duration
}

// New creates trm.Settings.
func New(oo ...Opt) (Settings, error) {
	s := &Settings{
		ctxKey:       nil,
		propagation:  nil,
		isCancelable: nil,
		timeout:      nil,
	}

	for _, o := range oo {
		if err := o(s); err != nil {
			return Settings{}, err
		}
	}

	return *s, nil
}

// Must returns Settings if err is nil and panics otherwise.
func Must(oo ...Opt) Settings {
	s, err := New(oo...)
	if err != nil {
		panic(err)
	}

	return s
}

// EnrichBy fills nil properties from external Settings.
func (s Settings) EnrichBy(external trm.Settings) (res trm.Settings) {
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

// CtxKey returns key(trm.CtxKey) to get trm.Transaction from context.Context.
func (s Settings) CtxKey() trm.CtxKey {
	if s.ctxKey == nil {
		return DefaultCtxKey
	}

	return *s.ctxKey
}

// CtxKeyOrNil returns key(*trm.CtxKey) or nil.
func (s Settings) CtxKeyOrNil() *trm.CtxKey {
	return s.ctxKey
}

// SetCtxKey sets key(*trm.CtxKey) to store trm.Transaction in context.Context.
func (s Settings) SetCtxKey(key *trm.CtxKey) trm.Settings {
	return s.setCtxKey(key)
}

func (s Settings) setCtxKey(key *trm.CtxKey) Settings {
	s.ctxKey = key

	return s
}

// Propagation returns trm.Propagation.
func (s Settings) Propagation() trm.Propagation {
	if s.propagation == nil {
		return defaultPropagation
	}

	return *s.propagation
}

// PropagationOrNil returns trm.Propagation or nil.
func (s Settings) PropagationOrNil() *trm.Propagation {
	return s.propagation
}

// SetPropagation sets *trm.Propagation.
func (s Settings) SetPropagation(p *trm.Propagation) trm.Settings {
	return s.setPropagation(p)
}

func (s Settings) setPropagation(p *trm.Propagation) Settings {
	s.propagation = p

	return s
}

// Cancelable defines that parent trm.Transaction can cancel children transactions.
func (s Settings) Cancelable() bool {
	if s.isCancelable == nil {
		return defaultCancelable
	}

	return *s.isCancelable
}

// CancelableOrNil returns Cancelable or nil.
func (s Settings) CancelableOrNil() *bool {
	return s.isCancelable
}

// SetCancelable sets Cancelable(*bool).
func (s Settings) SetCancelable(t *bool) trm.Settings {
	return s.setCancelable(t)
}

func (s Settings) setCancelable(t *bool) Settings {
	s.isCancelable = t

	return s
}

// TimeoutOrNil returns time.Duration of the trm.Transaction.
func (s Settings) TimeoutOrNil() *time.Duration {
	return s.timeout
}

// SetTimeout sets *time.Duration for trm.Transaction.
func (s Settings) SetTimeout(t *time.Duration) trm.Settings {
	return s.setTimeout(t)
}

func (s Settings) setTimeout(t *time.Duration) Settings {
	s.timeout = t

	return s
}
