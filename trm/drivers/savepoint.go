package drivers

import (
	"strconv"
	"sync/atomic"
)

const savePointPrefix = "tx_"

// SavePointCounter manages nested savepoint state atomically.
type SavePointCounter struct {
	saves int64
}

// NewSavePointCounter creates a SavePointCounter.
func NewSavePointCounter() *SavePointCounter {
	return &SavePointCounter{saves: 0}
}

// HasSavePoint reports whether there are active nested savepoints.
func (s *SavePointCounter) HasSavePoint() bool {
	return atomic.LoadInt64(&s.saves) > 0
}

// IncrementID atomically increments and returns the new savepoint identifier.
func (s *SavePointCounter) IncrementID() string {
	id := atomic.AddInt64(&s.saves, 1)

	return savePointPrefix + strconv.FormatInt(id, 10)
}

// DecrementID atomically decrements the counter and returns the identifier of the savepoint being released.
func (s *SavePointCounter) DecrementID() string {
	id := atomic.AddInt64(&s.saves, -1) + 1

	return savePointPrefix + strconv.FormatInt(id, 10)
}
