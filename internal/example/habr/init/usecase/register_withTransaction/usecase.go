package register

import (
	"github.com/jmoiron/sqlx"

	"init/domain"
	"init/queue"
	"init/repo"
)

type In struct {
	Username string
	Password string
}

type Register struct {
	userRepo domain.UserRepo
	db       *sqlx.DB
	queue    queue.Queue[domain.Registered]
}

func New(userRepo domain.UserRepo, db *sqlx.DB, queue queue.Queue[domain.Registered]) *Register {
	return &Register{userRepo: userRepo, db: db, queue: queue}
}

func (u *Register) Handle(tx *sqlx.Tx, in In) (user *domain.User, err error) {
	user, err = domain.NewUser(in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	err = repo.WithTransaction(tx, func(tx *sqlx.Tx) error {
		if err = u.userRepo.Save(tx, user); err != nil {
			return err
		}

		return u.queue.Publish(tx, domain.Registered{ID: user.ID})
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
