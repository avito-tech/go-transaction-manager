//go:build go1.14
// +build go1.14

package test

import (
	"testing"
)

// Cleanup is a helper function to register cleanup function for a test.
// t.Cleanup was added in go1.14.
func Cleanup(t *testing.T, fn func()) {
	t.Cleanup(fn)
}
