package domain

type UserRepo interface {
	GetByID(Tr, UserID) (*User, error)
	Save(Tr, *User) error
}

type OrderRepo interface {
	GetByID(Tr, OrderID) (*Order, error)
	GetByUserID(Tr, UserID) (*Order, error)
	Save(Tr, *Order) error
}
