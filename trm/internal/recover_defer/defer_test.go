package recover_defer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeferReturnErr(t *testing.T) {
	t.Parallel()

	err := func() (err error) {
		defer func() { err = errTest }()

		return nil
	}()

	assert.Error(t, err)
}

func TestDeferNotReturnErr(t *testing.T) {
	t.Parallel()

	//nolint:unparam
	err := func() error {
		defer func() error { return errTest }()

		return nil
	}()

	assert.NoError(t, err)
}
