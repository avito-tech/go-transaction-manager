package transaction

import "context"

// ErrBegin occurs when a transaction started with an error.
var ErrBegin = errTransaction("begin")

// Manager manages a transaction from Begin to Commit or Rollback.
type Manager interface {
	// Do processes a transaction inside a closure.
	Do(context.Context, func(ctx context.Context) error) error
	// DoWithSettings processes a transaction inside a closure with custom transaction.Settings.
	DoWithSettings(context.Context, Settings, func(ctx context.Context) error) error
}
