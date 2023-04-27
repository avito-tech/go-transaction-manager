package domain

import (
	"context"
)

type UserRepo interface {
	GetByID(context.Context, UserID) (*User, error)
	Save(context.Context, *User) error
}

type OrderRepo interface {
	GetByID(context.Context, OrderID) (*Order, error)
	GetByUserID(context.Context, UserID) (*Order, error)
	Save(context.Context, *Order) error
}
