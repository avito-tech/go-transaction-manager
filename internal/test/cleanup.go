//go:build !go1.14
// +build !go1.14

package test

import (
	"testing"
)

// Cleanup skipped for go1.13 and lower.
func Cleanup(t *testing.T, fn func()) {
}
