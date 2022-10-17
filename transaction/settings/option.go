package settings

import "github.com/avito-tech/go-transaction-manager/transaction"

// WithCtxKey sets up transaction.CtxKey for the transaction.Settings.
func WithCtxKey(key transaction.CtxKey) Opt {
	return func(s *settings) {
		s.ctxKey = &key
	}
}

// WithReadOnly sets up block to write to a database for the transaction.Settings.
func WithReadOnly(is bool) Opt {
	return func(s *settings) {
		s.isReadOnly = &is
	}
}

// WithPropagation sets up a transaction.Propagation for the transaction.Settings.
func WithPropagation(p transaction.Propagation) Opt {
	return func(s *settings) {
		s.propagation = &p
	}
}
