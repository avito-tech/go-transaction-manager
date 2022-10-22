//nolint:ireturn
package manager

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

type s struct{}

func (s s) EnrichBy(_ transaction.Settings) transaction.Settings {
	panic("implement me")
}

func (s s) CtxKey() transaction.CtxKey {
	panic("implement me")
}

func (s s) CtxKeyOrNil() *transaction.CtxKey {
	panic("implement me")
}

func (s s) SetCtxKey(_ *transaction.CtxKey) transaction.Settings {
	panic("implement me")
}

func (s s) IsReadOnly() bool {
	panic("implement me")
}

func (s s) IsReadOnlyOrNil() *bool {
	panic("implement me")
}

func (s s) SetIsReadOnly(b *bool) transaction.Settings {
	panic("implement me")
}

func (s s) Propagation() transaction.Propagation {
	panic("implement me")
}

func (s s) PropagationOrNil() *transaction.Propagation {
	panic("implement me")
}

func (s s) SetPropagation(_ *transaction.Propagation) transaction.Settings {
	panic("implement me")
}

func (s s) Cancelable() bool {
	panic("implement me")
}

func (s s) CancelableOrNil() *bool {
	panic("implement me")
}

func (s s) SetCancelable(_ *bool) transaction.Settings {
	panic("implement me")
}

func (s s) Timeout() time.Duration {
	panic("implement me")
}

func (s s) TimeoutOrNil() *time.Duration {
	panic("implement me")
}

func (s s) SetTimeout(_ *time.Duration) transaction.Settings {
	panic("implement me")
}
