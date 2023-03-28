package trm

import "context"

type Cmd func(ctx context.Context) error

// UoW is an implementation of Unit of Work.
type UoW interface {
	Register(ctx context.Context, cmd Cmd) error
	Commit(ctx context.Context) error
}
