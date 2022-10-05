// Package settings implements transaction.Settings.
package settings

import "github.com/avito-tech/go-transaction-manager/transaction"

// WithCtxKey sets up transaction.CtxKey for the Settings.
func WithCtxKey(key transaction.CtxKey) Opt {
	return func(s *Settings) {
		s.ctxKey = key
	}
}

// WithReadOnly sets up block to write to a database for the Settings.
func WithReadOnly(is bool) Opt {
	return func(s *Settings) {
		s.isReadOnly = is
	}
}

// WithPropagation sets up a transaction.Propagation for the Settings.
func WithPropagation(p transaction.Propagation) Opt {
	return func(s *Settings) {
		s.propagation = p
	}
}
