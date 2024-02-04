//go:build !go1.14
// +build !go1.14

package test

func Cleanup(t *testing.T, fn func()) {
}
