package context

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	const iterations = 50

	wg := sync.WaitGroup{}
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()

			Generate()
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(iterations+1), Generate().(int64))
}
