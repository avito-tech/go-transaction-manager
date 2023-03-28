//nolint:ireturn
package manager

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/trm"
)

type s struct{}

func (s s) EnrichBy(_ trm.Settings) trm.Settings {
	panic("implement me")
}

func (s s) CtxKey() trm.CtxKey {
	panic("implement me")
}

func (s s) CtxKeyOrNil() *trm.CtxKey {
	panic("implement me")
}

func (s s) SetCtxKey(_ *trm.CtxKey) trm.Settings {
	panic("implement me")
}

func (s s) IsReadOnly() bool {
	panic("implement me")
}

func (s s) IsReadOnlyOrNil() *bool {
	panic("implement me")
}

func (s s) SetIsReadOnly(_ *bool) trm.Settings {
	panic("implement me")
}

func (s s) Propagation() trm.Propagation {
	panic("implement me")
}

func (s s) PropagationOrNil() *trm.Propagation {
	panic("implement me")
}

func (s s) SetPropagation(_ *trm.Propagation) trm.Settings {
	panic("implement me")
}

func (s s) Cancelable() bool {
	panic("implement me")
}

func (s s) CancelableOrNil() *bool {
	panic("implement me")
}

func (s s) SetCancelable(_ *bool) trm.Settings {
	panic("implement me")
}

func (s s) Timeout() time.Duration {
	panic("implement me")
}

func (s s) TimeoutOrNil() *time.Duration {
	panic("implement me")
}

func (s s) SetTimeout(_ *time.Duration) trm.Settings {
	panic("implement me")
}
