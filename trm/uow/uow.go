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

func (u *uow) Register(_ context.Context, cmd trm.Cmd) (interface{}, error) {
	u.cmds = append(u.cmds, cmd)

	//nolint:nilnil
	return nil, nil
}

func (u *uow) Commit(ctx context.Context) ([]interface{}, error) {
	res := make([]interface{}, 0, len(u.cmds))

	err := u.manager.Do(ctx, func(ctx context.Context) error {
		for _, cmd := range u.cmds {
			item, err := cmd(ctx)
			if err != nil {
				return err
			}

			res = append(res, item)
		}

		u.cmds = nil

		return nil
	})

	return res, err
}
