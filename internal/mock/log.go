// Package mock implements dependencies for testing.
//
//revive:disable:unexported-return
package mock

import (
	"context"
)

type log struct {
	Logged []string
}

// NewLog create mock log.
func NewLog() *log {
	return &log{}
}

func (l *log) Warning(_ context.Context, msg string) {
	l.Logged = append(l.Logged, msg)
}

type zeroLog struct{}

// NewZeroLog create mock log, which should not be called.
func NewZeroLog() *zeroLog {
	return &zeroLog{}
}

func (l *zeroLog) Warning(_ context.Context, _ string) {
	panic("logger should not be called")
}
