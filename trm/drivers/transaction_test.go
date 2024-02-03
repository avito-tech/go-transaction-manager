package drivers

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsClose(t *testing.T) {
	errExpected := errors.New("expected")
	err := errors.New("test")

	isClosed := NewIsClosed()

	require.True(t, isClosed.IsActive())
	require.False(t, isClosed.IsClosed())

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			require.NotPanics(t, func() {
				<-isClosed.Closed()

				require.ErrorIs(t, isClosed.Err(), errExpected)
				require.False(t, isClosed.IsActive())
				require.True(t, isClosed.IsClosed())

				isClosed.CloseWithCause(err)

				require.ErrorIs(t, isClosed.Err(), errExpected)
				require.False(t, isClosed.IsActive())
				require.True(t, isClosed.IsClosed())
			})
		}()
	}

	isClosed.CloseWithCause(errExpected)
	isClosed.Close()

	wg.Wait()

	require.ErrorIs(t, isClosed.Err(), errExpected)
	require.False(t, isClosed.IsActive())
	require.True(t, isClosed.IsClosed())
}

func TestIsClose_Err_nil(t *testing.T) {
	isClosed := NewIsClosed()

	isClosed.Close()

	require.NoError(t, isClosed.Err())
}
