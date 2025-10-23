// Package drivers contains instruments for drivers.
package drivers

import (
	"errors"
	"sync"
)

// ErrRollbackTr is used to rollback transaction, which runs in closure.
var ErrRollbackTr = errors.New("rollback")

// IsClosed stores the state of the trm.Transaction.
type IsClosed struct {
	isClosed bool
	mu       sync.RWMutex
	once     sync.Once
	ch       chan struct{}
	// txErr stores error from transaction commit or rollback
	err error
}

// NewIsClosed creates a new IsClosed.
func NewIsClosed() *IsClosed {
	return &IsClosed{
		isClosed: false,
		mu:       sync.RWMutex{},
		once:     sync.Once{},
		ch:       make(chan struct{}),
		err:      nil,
	}
}

// IsActive returns true if the channel is open.
func (a *IsClosed) IsActive() bool {
	return !a.IsClosed()
}

// IsClosed returns true if the channel is closed.
func (a *IsClosed) IsClosed() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.isClosed
}

// Closed returns a channel that's closed when trm.Transaction done.
// Closed is provided for use in select statements.
func (a *IsClosed) Closed() <-chan struct{} {
	return a.ch
}

// Close closes the channel.
func (a *IsClosed) Close() {
	a.CloseWithCause(nil)
}

// CloseWithCause closes the channel and stores error from transaction commit or rollback.
func (a *IsClosed) CloseWithCause(err error) {
	a.once.Do(func() {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.isClosed = true
		a.err = err

		close(a.ch)
	})
}

// Err is inspired function Err in https://pkg.go.dev/context#Context
func (a *IsClosed) Err() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.err
}
