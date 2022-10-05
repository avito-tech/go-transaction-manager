package transaction

import "context"

// ErrBegin occurs when a transaction started with an error.
var ErrBegin = errTransaction("begin")

// Manager manages a transaction from Begin to Commit or Rollback.
type Manager interface {
	// Do processes a transaction inside a closure.
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
