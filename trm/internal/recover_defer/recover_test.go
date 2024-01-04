package recover_defer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const panicText = "panic"

var errTest = errors.New("err")

func TestInline(t *testing.T) {
	t.Parallel()

	defer func() {
		assert.Equal(t, panicText, recover())
	}()

	panic(panicText)
}

func recoverer(t *testing.T) func() {
	return func() {
		assert.Equal(t, panicText, recover())
	}
}

func TestNestedFunc(t *testing.T) {
	t.Parallel()

	d := recoverer(t)

	defer d()

	panic("panic")
}

func emptyRecoverer(t *testing.T) func() {
	return func() {
		assert.Empty(t, recover())
	}
}

func TestNestedFuncInline(t *testing.T) {
	t.Parallel()

	d := emptyRecoverer(t)

	defer func() {
		d()

		assert.NotNil(t, recover())
	}()

	panic("panic")
}
