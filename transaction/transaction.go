// Package transaction is an interface to create a transactional usecase
// in the Application layer.
//
//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package transaction

import (
	"context"
	"errors"
)

// ErrTransaction is an error while working with a transaction.
var ErrTransaction = errors.New("transaction error")

// Manager manages a transaction from Begin to Commit or Rollback.
type Manager interface {
	// Do starts a transaction inside a closure.
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
