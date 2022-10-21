package manager

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/transaction"
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

// NewChained creates *ChainedMW or chained transaction.Manager.
func NewChained(mm []transaction.Manager) *ChainedMW {
	if len(mm) == 0 {
		return &ChainedMW{
			do:             nilNextDo,
			doWithSettings: nilNextDoWithSettings,
		}
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
	}
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

func newNextDo(m transaction.Manager, n nextDo) nextDo {
	return func(ctx context.Context, fn callback) error {
		return m.Do(ctx, func(ctx context.Context) error {
			return n(ctx, fn)
		})
	}
}

func newLastDo(m transaction.Manager) nextDo {
	return func(ctx context.Context, fn callback) error {
		return m.Do(ctx, fn)
	}
}

// DoWithSettings is an implementation of transaction.Manager.
//
// WARNING: transaction.CtxKey should not be set in transaction.Settings otherwise all transaction.Manager would get same transaction.Transaction from context.Context.
func (c *ChainedMW) DoWithSettings(
	ctx context.Context,
	s transaction.Settings,
	fn func(ctx context.Context) error,
) error {
	return c.doWithSettings(ctx, s, fn)
}

type nextDoWithSettings func(context.Context, transaction.Settings, callback) error

func nilNextDoWithSettings(
	ctx context.Context,
	_ transaction.Settings,
	fn callback,
) error {
	return fn(ctx)
}

func newNextDoWithSettings(m transaction.Manager, n nextDoWithSettings) nextDoWithSettings {
	return func(
		ctx context.Context,
		s transaction.Settings,
		fn callback,
	) error {
		return m.DoWithSettings(ctx, s, func(ctx context.Context) error {
			return n(ctx, s, fn)
		})
	}
}

func newLastDoWithSettings(m transaction.Manager) nextDoWithSettings {
	return func(
		ctx context.Context,
		s transaction.Settings,
		fn callback,
	) error {
		return m.DoWithSettings(ctx, s, fn)
	}
}
