package drivers

import (
	"sync"
	"sync/atomic"
)

type IsActive struct {
	isActive int64
	ch       chan struct{}
	close    func()
}

func NewIsActive() *IsActive {
	once := sync.Once{}
	closeCh := make(chan struct{})

	return &IsActive{
		isActive: 1,
		ch:       closeCh,
		close: func() {
			once.Do(func() {
				close(closeCh)
			})
		}}
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (a *IsActive) IsActive() bool {
	return atomic.LoadInt64(&a.isActive) == 1
}

// Deactivated returns a channel that's closed when trm.Transaction done.
// Deactivated is provided for use in select statements
func (a *IsActive) Deactivated() <-chan struct{} {
	return a.ch
}

func (a *IsActive) Deactivate() {
	a.close()

	atomic.SwapInt64(&a.isActive, 0)
}
