package purchase

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"

	"finish/domain"
	"finish/queue"
)

type In struct {
	UserID    domain.UserID
	ProductID domain.ProductID
	Quantity  int64
}

type Purchase struct {
	orderRepo domain.OrderRepo
	trm       trm.Manager
	queue     queue.Queue[domain.Purchased]
}

func New(orderRepo domain.OrderRepo, trm trm.Manager, queue queue.Queue[domain.Purchased]) *Purchase {
	return &Purchase{orderRepo: orderRepo, trm: trm, queue: queue}
}

func (u *Purchase) Handle(ctx context.Context, in In) (order *domain.Order, err error) {
	order, err = domain.NewOrder(in.UserID, in.ProductID, in.Quantity)
	if err != nil {
		return nil, err
	}

	err = u.trm.Do(ctx, func(ctx context.Context) error {
		if err = u.orderRepo.Save(ctx, order); err != nil {
			return err
		}

		return u.queue.Publish(ctx, domain.Purchased{ID: order.ID})
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}
