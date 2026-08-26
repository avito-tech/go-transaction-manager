package drivers

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePointCounter(t *testing.T) {
	t.Parallel()

	sp := NewSavePointCounter()
	assert.False(t, sp.HasSavePoint())

	assert.Equal(t, "tx_1", sp.IncrementID())
	assert.True(t, sp.HasSavePoint())

	assert.Equal(t, "tx_2", sp.IncrementID())
	assert.True(t, sp.HasSavePoint())

	assert.Equal(t, "tx_2", sp.DecrementID())
	assert.True(t, sp.HasSavePoint())

	assert.Equal(t, "tx_1", sp.DecrementID())
	assert.False(t, sp.HasSavePoint())
}

func TestSavePointCounter_Concurrency(t *testing.T) {
	t.Parallel()

	sp := NewSavePointCounter()
	const iterations = 500
	generated := sync.Map{}
	wg := sync.WaitGroup{}
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()

			id := sp.IncrementID()
			_, loaded := generated.LoadOrStore(id, struct{}{})
			assert.False(t, loaded, "duplicate savepoint generated: %s", id)
		}()
	}

	wg.Wait()
	assert.True(t, sp.HasSavePoint())
	assert.Equal(t, int64(iterations), sp.saves)
}
