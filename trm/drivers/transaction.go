package drivers

import (
	"errors"
	"sync"
)

// ErrRollbackTr is used to rollback transaction, which runs in closure.
var ErrRollbackTr = errors.New("rollback")

type IsClose struct {
	isClosed bool
	mu       sync.RWMutex
	once     sync.Once
	ch       chan struct{}
	// txErr stores error from transaction commit or rollback
	err error
}

func NewIsClosed() *IsClose {
	return &IsClose{
		isClosed: false,
		mu:       sync.RWMutex{},
		once:     sync.Once{},
		ch:       make(chan struct{}),
	}
}

func (a *IsClose) IsActive() bool {
	return !a.IsClosed()
}

func (a *IsClose) IsClosed() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.isClosed
}

// Deactivated returns a channel that's closed when trm.Transaction done.
// Deactivated is provided for use in select statements
func (a *IsClose) Closed() <-chan struct{} {
	return a.ch
}

func (a *IsClose) Close() {
	a.CloseWithCause(nil)
}

func (a *IsClose) CloseWithCause(err error) {
	a.once.Do(func() {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.isClosed = true
		a.err = err

		close(a.ch)
	})
}

func (a *IsClose) Err() error {
	a.mu.RLock()
	defer a.mu.RLock()

	return a.err
}
