package fast_purchase

import (
	"github.com/jmoiron/sqlx"

	"init/domain"
	"init/usecase/purchase"
	"init/usecase/register"
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
	db       *sqlx.DB
	register *register.Register
	purchase *purchase.Purchase
}

func New(db *sqlx.DB, register *register.Register, purchase *purchase.Purchase) *FastPurchase {
	return &FastPurchase{db: db, register: register, purchase: purchase}
}

func (u *FastPurchase) Handle(in In) (out Out, err error) {
	tx, err := u.db.Beginx()
	if err != nil {
		return Out{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if out.User, err = u.register.Handle(tx, register.In{
		Username: in.Register.Username,
	}); err != nil {
		return Out{}, err
	}

	if out.Order, err = u.purchase.Handle(tx, purchase.In{
		UserID:    out.User.ID,
		ProductID: in.Purchase.ProductID,
		Quantity:  in.Purchase.Quantity,
	}); err != nil {
		return Out{}, err
	}

	if err = tx.Commit(); err != nil {
		return Out{}, err
	}

	return out, nil
}
