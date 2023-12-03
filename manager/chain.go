package manager

import (
	"context"

	trm "github.com/avito-tech/go-transaction-manager/v2"
)

// ChainedMW starts transactions in the order given and commit/rollback in reverse order.
//
// WARNING: Rollback of last transactions isn't affected done commits.
// ChainedMW should be only used if the application can tolerate or
// recover from an inconsistent state caused by partially committed transactions.
type ChainedMW struct {
	do             nextDo
	doWithSettings nextDoWithSettings
}

// NewChained creates *ChainedMW or chained trm.Manager.
func NewChained(mm []trm.Manager, _ ...Opt) (*ChainedMW, error) {
	if len(mm) == 0 {
		return &ChainedMW{
			do:             nilNextDo,
			doWithSettings: nilNextDoWithSettings,
		}, nil
	}

	last := len(mm) - 1
	do := newLastDo(mm[last])
	doWithSettings := newLastDoWithSettings(mm[last])

	for index := last - 1; index >= 0; index-- {
		do = newNextDo(mm[index], do)
		doWithSettings = newNextDoWithSettings(mm[index], doWithSettings)
	}

	return &ChainedMW{
		do:             do,
		doWithSettings: doWithSettings,
	}, nil
}

// MustChained returns ChainedMW if err is nil and panics otherwise.
func MustChained(mm []trm.Manager, oo ...Opt) *ChainedMW {
	s, err := NewChained(mm, oo...)
	if err != nil {
		panic(err)
	}

	return s
}

//revive:disable:exported
func (c *ChainedMW) Do(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return c.do(ctx, fn)
}

type callback func(ctx context.Context) error

type nextDo func(ctx context.Context, fn callback) error

func nilNextDo(ctx context.Context, fn callback) error {
	return fn(ctx)
}

func newNextDo(m trm.Manager, n nextDo) nextDo {
	return func(ctx context.Context, fn callback) error {
		return m.Do(ctx, func(ctx context.Context) error {
			return n(ctx, fn)
		})
	}
}

func newLastDo(m trm.Manager) nextDo {
	return func(ctx context.Context, fn callback) error {
		return m.Do(ctx, fn)
	}
}

// DoWithSettings is an implementation of trm.Manager.
//
// WARNING: trm.CtxKey should not be set in trm.Settings otherwise all trm.Manager would get same trm.Transaction from context.Context.
func (c *ChainedMW) DoWithSettings(
	ctx context.Context,
	s trm.Settings,
	fn func(ctx context.Context) error,
) error {
	return c.doWithSettings(ctx, s, fn)
}

type nextDoWithSettings func(context.Context, trm.Settings, callback) error

func nilNextDoWithSettings(
	ctx context.Context,
	_ trm.Settings,
	fn callback,
) error {
	return fn(ctx)
}

func newNextDoWithSettings(m trm.Manager, n nextDoWithSettings) nextDoWithSettings {
	return func(
		ctx context.Context,
		s trm.Settings,
		fn callback,
	) error {
		return m.DoWithSettings(ctx, s, func(ctx context.Context) error {
			return n(ctx, s, fn)
		})
	}
}

func newLastDoWithSettings(m trm.Manager) nextDoWithSettings {
	return func(
		ctx context.Context,
		s trm.Settings,
		fn callback,
	) error {
		return m.DoWithSettings(ctx, s, fn)
	}
}
