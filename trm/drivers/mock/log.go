// Package mock implements dependencies for testing.
// Deprecated: You should NOT use this package in your application code.
//
//revive:disable:unexported-return
//revive:disable:exported
package mock

import (
	"context"
)

type Log struct {
	Logged []string
}

// NewLog create mock Log.
func NewLog() *Log {
	return &Log{
		Logged: nil,
	}
}

func (l *Log) Warning(_ context.Context, msg string) {
	l.Logged = append(l.Logged, msg)
}

type zeroLog struct{}

// NewZeroLog create mock Log, which should not be called.
func NewZeroLog() *zeroLog {
	return &zeroLog{}
}

func (l *zeroLog) Warning(_ context.Context, _ string) {
	panic("logger should not be called")
}
