package transaction

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import "context"

// Manager manages a transaction from Begin to Commit or Rollback.
type Manager interface {
	// Do processes a transaction inside a closure.
	Do(context.Context, func(ctx context.Context) error) error
	// DoWithSettings processes a transaction inside a closure with custom transaction.Settings.
	DoWithSettings(context.Context, Settings, func(ctx context.Context) error) error
}
