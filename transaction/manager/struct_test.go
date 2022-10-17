//nolint:ireturn
package manager

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/transaction"
)

type l struct{}

func (l l) Printf(format string, a ...interface{}) {
	panic("implement me")
}

type s struct{}

func (s s) EnrichBy(external transaction.Settings) transaction.Settings {
	panic("implement me")
}

func (s s) CtxKey() transaction.CtxKey {
	panic("implement me")
}

func (s s) CtxKeyOrNil() *transaction.CtxKey {
	panic("implement me")
}

func (s s) SetCtxKey(key *transaction.CtxKey) transaction.Settings {
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

func (s s) SetPropagation(propagation *transaction.Propagation) transaction.Settings {
	panic("implement me")
}

func (s s) Timeout() time.Duration {
	panic("implement me")
}

func (s s) TimeoutOrNil() *time.Duration {
	panic("implement me")
}

func (s s) SetTimeout(duration *time.Duration) transaction.Settings {
	panic("implement me")
}
