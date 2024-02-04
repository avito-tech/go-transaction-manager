//go:build go1.14
// +build go1.14

package test

import (
	"testing"
)

func Cleanup(t *testing.T, fn func()) {
	t.Cleanup(fn)
}
