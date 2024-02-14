//go:build go1.20
// +build go1.20

package goredis8

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	// https://github.com/redis/go-redis/issues/1029
	goleak.VerifyTestMain(m, goleak.IgnoreAnyFunction("github.com/go-redis/redis/v8/internal/pool.(*ConnPool).reaper"))
}
