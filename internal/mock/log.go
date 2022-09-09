// Package mock implements dependencies for testing.
//
//revive:disable:unexported-return
package mock

import "fmt"

type log struct {
	Logged []string
}

// NewLog create mock log.
func NewLog() *log {
	return &log{}
}

func (l *log) Printf(format string, a ...interface{}) {
	l.Logged = append(l.Logged, fmt.Sprintf(format, a...))
}

type zeroLog struct{}

// NewZeroLog create mock log, which should not be called.
func NewZeroLog() *zeroLog {
	return &zeroLog{}
}

func (l *zeroLog) Printf(_ string, _ ...interface{}) {
	panic("logger should not be called")
}
