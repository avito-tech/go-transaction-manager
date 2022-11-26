package redis

import "github.com/go-redis/redis/v8"

// WithMulti sets up sql.TxOptions for the Settings.
func WithMulti(isMulti bool) Opt {
	return func(s *Settings) error {
		*s = s.setIsMulti(&isMulti)

		return nil
	}
}

// WithWatchKeys sets WATCH keys in Transaction.
func WithWatchKeys(keys ...string) Opt {
	return func(s *Settings) error {
		*s = s.setWatchKeys(keys)

		return nil
	}
}

// WithTxDecorator sets TxDecorator to change behavior of Transaction.
func WithTxDecorator(in TxDecorator) Opt {
	return func(s *Settings) error {
		*s = s.setTxDecorator(in)

		return nil
	}
}

// WithRet sets link on []redis.Cmder to get responses of commands in Transaction
// WARNING: Responses don't clean automatically, use WithRet only with DoWithSettings of trm.Manager.
func WithRet(in *[]redis.Cmder) Opt {
	return func(s *Settings) error {
		*s = s.SetReturn(in)

		return nil
	}
}
