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

func TestKeyGenerator_Generate_Concurrency(t *testing.T) {
	t.Parallel()

	const iterations = 1000
	generator := NewKeyGenerator()
	generatedKeys := sync.Map{}

	wg := sync.WaitGroup{}
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()

			key := generator.Generate()
			_, loaded := generatedKeys.LoadOrStore(key, struct{}{})
			assert.False(t, loaded, "duplicate key generated: %v", key)
		}()
	}

	wg.Wait()
}
