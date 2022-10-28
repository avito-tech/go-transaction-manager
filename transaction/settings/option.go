package settings

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

// WithCtxKey sets up transaction.CtxKey for the transaction.Settings.
func WithCtxKey(key transaction.CtxKey) Opt {
	return func(s *Settings) {
		*s = s.setCtxKey(&key)
	}
}

// WithPropagation sets up a transaction.Propagation for the transaction.Settings.
func WithPropagation(p transaction.Propagation) Opt {
	return func(s *Settings) {
		*s = s.setPropagation(&p)
	}
}

// WithCancelable sets up possibility to cancel child goroutines when parent transaction.Transaction was canceled.
func WithCancelable(t bool) Opt {
	return func(s *Settings) {
		*s = s.setCancelable(&t)
	}
}

// WithTimeout sets up a timeout for the transaction.Transaction.
func WithTimeout(t time.Duration) Opt {
	return func(s *Settings) {
		*s = s.setTimeout(&t)
	}
}
