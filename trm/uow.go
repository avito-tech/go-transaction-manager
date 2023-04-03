package trm

import "context"

type Cmd func(ctx context.Context) (interface{}, error)

// UoW is an implementation of Unit of Work.
type UoW interface {
	Register(ctx context.Context, cmd Cmd) (interface{}, error)
	Commit(ctx context.Context) ([]interface{}, error)
}
