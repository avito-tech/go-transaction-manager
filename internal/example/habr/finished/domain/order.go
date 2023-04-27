package domain

import "errors"

var (
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidProductID = errors.New("invalid product id")
	ErrInvalidQuantity  = errors.New("quantity should be more 1")
)

type OrderID int64
type ProductID int64

type Order struct {
	ID        OrderID
	ProductID ProductID
	UserID    UserID
	Quantity  int64
}

func NewOrder(userID UserID, productID ProductID, quantity int64) (*Order, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	if productID <= 0 {
		return nil, ErrInvalidProductID
	}

	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	return &Order{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}, nil
}

type Purchased struct {
	ID OrderID
}
