package uow

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
)

func NewUoW(manager trm.Manager) *uow {
	return &uow{manager: manager}
}

type uow struct {
	manager trm.Manager
	cmds    []trm.Cmd
}

func (u *uow) Register(_ context.Context, cmd trm.Cmd) error {
	u.cmds = append(u.cmds, cmd)

	return nil
}

func (u *uow) Commit(ctx context.Context) error {
	return u.manager.Do(ctx, func(ctx context.Context) error {
		for _, cmd := range u.cmds {
			if err := cmd(ctx); err != nil {
				return err
			}
		}

		u.cmds = nil

		return nil
	})
}
