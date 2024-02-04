// Package trm contains of interfaces to programmatic transaction management.
package trm

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import (
	"context"
	"errors"

	"go.uber.org/multierr"
)

// Manager manages a transaction from Begin to Commit or Rollback.
type Manager interface {
	// Do processes a transaction inside a closure.
	Do(context.Context, func(ctx context.Context) error) error
	// DoWithSettings processes a transaction inside a closure with custom trm.Settings.
	DoWithSettings(context.Context, Settings, func(ctx context.Context) error) error
}

// ErrSkip marks error to skip rollback for transaction because of inside error.
var ErrSkip = errors.New("skippable")

// Skippable marks error as ErrSkip.
func Skippable(err error) error {
	if err == nil {
		return nil
	}

	return multierr.Append(err, ErrSkip)
}

// UnSkippable removes ErrSkip from error.
func UnSkippable(err error) error {
	if err == nil || !IsSkippable(err) {
		return err
	}

	ee := multierr.Errors(err)
	res := make([]error, 0, len(ee))

	for _, e := range ee {
		//nolint:errorlint,goerr113
		if e != ErrSkip {
			res = append(res, e)
		}
	}

	return multierr.Combine(res...)
}

// IsSkippable checks that the error is ErrSkip.
func IsSkippable(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrSkip)
}
