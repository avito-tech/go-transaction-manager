package transaction

import "time"

// Settings is the configuration of the Manager.
type Settings interface {
	CtxKey() CtxKey
	// TODO
	IsReadOnly() bool
	// TODO
	Propagation() Propagation
	// TODO
	Timeout() time.Duration
}
