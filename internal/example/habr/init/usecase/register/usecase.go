package register

import (
	"github.com/jmoiron/sqlx"

	"init/domain"
	"init/queue"
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
	hasExternalTx := true
	tr := domain.Tr(tx)
	if tx == nil {
		if tr, err = u.db.Beginx(); err != nil {
			return nil, err
		}
		hasExternalTx = false
	}

	defer func() {
		if !hasExternalTx && err != nil {
			tx.Rollback()
		}
	}()

	user, err = domain.NewUser(in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	if err = u.userRepo.Save(tr, user); err != nil {
		return nil, err
	}

	if err = u.queue.Publish(tr, domain.Registered{ID: user.ID}); err != nil {
		return nil, err
	}

	if !hasExternalTx {
		return nil, tx.Commit()
	}

	return user, nil
}
