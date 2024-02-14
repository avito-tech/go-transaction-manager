package goredis8

import (
	"sync"

	"github.com/go-redis/redis/v8"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

const (
	// DefaultMulti is a default value for Settings.IsMulti.
	DefaultMulti = true
)

// Opt is a type to configure Settings.
type Opt func(*Settings) error

// Settings contains settings for redis.Transaction.
type Settings struct {
	trm.Settings
	isMulti     *bool
	watchKeys   []string
	txDecorator []TxDecorator

	ret   *[]redis.Cmder
	muRet sync.RWMutex
}

// NewSettings creates Settings.
func NewSettings(trms trm.Settings, oo ...Opt) (*Settings, error) {
	s := &Settings{
		Settings:    trms,
		isMulti:     nil,
		watchKeys:   nil,
		txDecorator: nil,
		ret:         nil,
		muRet:       sync.RWMutex{},
	}

	for _, o := range oo {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// MustSettings returns Settings if err is nil and panics otherwise.
func MustSettings(trms trm.Settings, oo ...Opt) *Settings {
	s, err := NewSettings(trms, oo...)
	if err != nil {
		panic(err)
	}

	return s
}

// EnrichBy fills nil properties from external Settings.
func (s *Settings) EnrichBy(in trm.Settings) trm.Settings {
	external, ok := in.(*Settings)
	if ok {
		if s.IsMultiOrNil() == nil {
			s = s.SetIsMulti(external.IsMultiOrNil())
		}

		if s.WatchKeys() == nil {
			s = s.SetWatchKeys(external.WatchKeys())
		}

		if s.TxDecorators() == nil {
			s = s.SetTxDecorators(external.TxDecorators()...)
		}

		if s.ReturnPtr() == nil {
			s = s.SetReturn(external.ReturnPtr())
		}
	}

	s.Settings = s.Settings.EnrichBy(in)

	return s
}

// IsMulti - true uses redis MULTI cmd.
func (s *Settings) IsMulti() bool {
	if s.isMulti == nil {
		return DefaultMulti
	}

	return *s.isMulti
}

// IsMultiOrNil returns IsMulti or nil.
func (s *Settings) IsMultiOrNil() *bool {
	return s.isMulti
}

// SetIsMulti set using or not Multi for transaction, see https://redis.uptrace.dev/guide/go-redis-pipelines.html#transactions.
func (s *Settings) SetIsMulti(in *bool) *Settings {
	return s.setIsMulti(in)
}

func (s *Settings) setIsMulti(in *bool) *Settings {
	s.isMulti = in

	return s
}

// WatchKeys returns keys for watching.
func (s *Settings) WatchKeys() []string {
	return s.watchKeys
}

// SetWatchKeys sets keys for watching, see https://redis.uptrace.dev/guide/go-redis-pipelines.html#watch.
func (s *Settings) SetWatchKeys(in []string) *Settings {
	return s.setWatchKeys(in)
}

func (s *Settings) setWatchKeys(in []string) *Settings {
	s.watchKeys = in

	return s
}

// TxDecorators returns TxDecorator decorators.
func (s *Settings) TxDecorators() []TxDecorator {
	return s.txDecorator
}

// SetTxDecorators sets TxDecorator decorators.
func (s *Settings) SetTxDecorators(in ...TxDecorator) *Settings {
	return s.setTxDecorator(in...)
}

func (s *Settings) setTxDecorator(in ...TxDecorator) *Settings {
	s.txDecorator = in

	return s
}

// ReturnPtr returns link to save []redis.Cmder from Transaction.
func (s *Settings) ReturnPtr() *[]redis.Cmder {
	s.muRet.RLock()
	defer s.muRet.RUnlock()

	return s.ret
}

// Return returns []redis.Cmder from Transaction.
func (s *Settings) Return() []redis.Cmder {
	res := s.ReturnPtr()
	if res != nil {
		return *s.ReturnPtr()
	}

	return nil
}

// AppendReturn append []redis.Cmder from Transaction.
func (s *Settings) AppendReturn(cmds ...redis.Cmder) {
	if s.ReturnPtr() == nil {
		return
	}

	s.muRet.Lock()
	defer s.muRet.Unlock()

	*s.ret = append(*s.ret, cmds...)
}

// SetReturn sets link to save []redis.Cmder from Transaction.
func (s *Settings) SetReturn(in *[]redis.Cmder) *Settings {
	return s.setReturn(in)
}

func (s *Settings) setReturn(in *[]redis.Cmder) *Settings {
	s.ret = in

	return s
}
