package fast_purchase

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"

	"finish/domain"
	"finish/usecase/purchase"
	"finish/usecase/register"
)

type In struct {
	Register RegisterIn
	Purchase PurchaseIn
}

type RegisterIn struct {
	Username string
	Password string
}

type PurchaseIn struct {
	ProductID domain.ProductID
	Quantity  int64
}

type Out struct {
	User  *domain.User
	Order *domain.Order
}

type FastPurchase struct {
	trm      trm.Manager
	register *register.Register
	purchase *purchase.Purchase
}

func New(trm trm.Manager, register *register.Register, purchase *purchase.Purchase) *FastPurchase {
	return &FastPurchase{trm: trm, register: register, purchase: purchase}
}

func (u *FastPurchase) Handle(ctx context.Context, in In) (out Out, err error) {
	err = u.trm.Do(ctx, func(ctx context.Context) error {
		out.User, err = u.register.Handle(ctx, register.In{
			Username: in.Register.Username,
			Password: in.Register.Password,
		})
		if err != nil {
			return err
		}

		out.Order, err = u.purchase.Handle(ctx, purchase.In{
			UserID:    out.User.ID,
			ProductID: in.Purchase.ProductID,
			Quantity:  in.Purchase.Quantity,
		})

		return err
	})
	if err != nil {
		return Out{}, err
	}

	return out, nil
}
