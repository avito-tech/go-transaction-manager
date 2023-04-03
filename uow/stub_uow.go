package uow

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
)

type stubUoW struct{}

func (u stubUoW) Register(ctx context.Context, cmd trm.Cmd) (interface{}, error) {
	return cmd(ctx)
}

func (u stubUoW) Commit(_ context.Context) ([]interface{}, error) {
	return nil, nil
}
