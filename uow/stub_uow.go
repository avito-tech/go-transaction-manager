package uow

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
)

type stubUoW struct{}

func (u stubUoW) Register(ctx context.Context, cmd trm.Cmd) error {
	return cmd(ctx)
}

func (u stubUoW) Commit(ctx context.Context) error {
	return nil
}
