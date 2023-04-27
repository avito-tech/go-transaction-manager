package domain

import "errors"

var (
	ErrEmptyUsername = errors.New("empty username")
	ErrEmptyPassword = errors.New("empty password")
)

type UserID int64

type User struct {
	ID       UserID
	Username string
	// В реальном проекте Пароли храните в захэшированном виде!
	Password     string
	Notification Notification
}

func NewUser(username string, password string) (*User, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}

	if password == "" {
		return nil, ErrEmptyPassword
	}

	return &User{
		Username: username,
		Password: password,
		Notification: Notification{
			Email: false,
			SMS:   true,
		},
	}, nil
}

type Notification struct {
	Email bool
	SMS   bool
}

type Registered struct {
	ID UserID
}
