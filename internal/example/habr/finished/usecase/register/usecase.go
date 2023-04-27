package register

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"

	"finish/domain"
	"finish/queue"
)

type In struct {
	Username string
	Password string
}

type Register struct {
	userRepo domain.UserRepo
	trm      trm.Manager
	queue    queue.Queue[domain.Registered]
}

func New(userRepo domain.UserRepo, trm trm.Manager, queue queue.Queue[domain.Registered]) *Register {
	return &Register{userRepo: userRepo, trm: trm, queue: queue}
}

func (u *Register) Handle(ctx context.Context, in In) (user *domain.User, err error) {
	user, err = domain.NewUser(in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	err = u.trm.Do(ctx, func(ctx context.Context) error {
		if err = u.userRepo.Save(ctx, user); err != nil {
			return err
		}

		return u.queue.Publish(ctx, domain.Registered{ID: user.ID})
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
