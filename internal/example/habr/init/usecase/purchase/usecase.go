package purchase

import (
	"github.com/jmoiron/sqlx"

	"init/domain"
	"init/queue"
)

type In struct {
	UserID    domain.UserID
	ProductID domain.ProductID
	Quantity  int64
}

type Purchase struct {
	orderRepo domain.OrderRepo
	db        *sqlx.DB
	queue     queue.Queue[domain.Purchased]
}

func New(orderRepo domain.OrderRepo, db *sqlx.DB, queue queue.Queue[domain.Purchased]) *Purchase {
	return &Purchase{orderRepo: orderRepo, db: db, queue: queue}
}

func (u *Purchase) Handle(tx *sqlx.Tx, in In) (order *domain.Order, err error) {
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

	order, err = domain.NewOrder(in.UserID, in.ProductID, in.Quantity)
	if err != nil {
		return nil, err
	}

	if err = u.orderRepo.Save(tr, order); err != nil {
		return nil, err
	}

	if err = u.queue.Publish(tr, domain.Purchased{ID: order.ID}); err != nil {
		return nil, err
	}

	if !hasExternalTx {
		return nil, tx.Commit()
	}

	return order, nil
}
