package settings

import (
	"time"

	trm "github.com/avito-tech/go-transaction-manager/v2"
)

// WithCtxKey sets up trm.CtxKey for the trm.Settings.
func WithCtxKey(key trm.CtxKey) Opt {
	return func(s *Settings) error {
		*s = s.setCtxKey(&key)

		return nil
	}
}

// WithPropagation sets up a trm.Propagation for the trm.Settings.
func WithPropagation(p trm.Propagation) Opt {
	return func(s *Settings) error {
		*s = s.setPropagation(&p)

		return nil
	}
}

// WithCancelable sets up possibility to cancel child goroutines when parent trm.Transaction was canceled.
func WithCancelable(t bool) Opt {
	return func(s *Settings) error {
		*s = s.setCancelable(&t)

		return nil
	}
}

// WithTimeout sets up a timeout for the trm.trm.
func WithTimeout(t time.Duration) Opt {
	return func(s *Settings) error {
		*s = s.setTimeout(&t)

		return nil
	}
}
