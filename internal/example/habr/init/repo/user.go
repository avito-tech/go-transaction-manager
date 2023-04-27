package repo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"init/domain"
)

type userRow struct {
	ID           domain.UserID `db:"id"`
	Username     string        `db:"username"`
	Password     string        `db:"password"`
	Notification notification  `db:"notification"`
}

type notification struct {
	Email bool `json:"email"`
	SMS   bool `json:"sms"`
}

func (n notification) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, &n)
	case string:
		return json.Unmarshal([]byte(v), &n)
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (n notification) Value() (driver.Value, error) {
	return json.Marshal(n)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetByID(tr domain.Tr, id domain.UserID) (*domain.User, error) {
	query := `SELECT * FROM "user" WHERE id = ?;`
	uRow := userRow{}

	if tr != nil {
		tr = r.db
	}

	if err := tr.Get(&uRow, r.db.Rebind(query), id); err != nil {
		return nil, err
	}

	return r.toModel(uRow), nil
}

func (r *userRepo) Save(tr domain.Tr, u *domain.User) error {
	query := `INSERT INTO "user" (username, password, notification)
VALUES (:username, :password, :notification)
ON CONFLICT (id)
    DO UPDATE SET username     = excluded.username,
                  password     = excluded.password,
                  notification = excluded.notification
RETURNING id;`

	if tr != nil {
		tr = r.db
	}

	rows, err := sqlx.NamedQuery(tr, query, r.toRow(u))
	if err != nil {
		return err
	}

	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}

	err = rows.Scan(&u.ID)

	return err
}

func (r userRepo) toModel(row userRow) *domain.User {
	return &domain.User{
		ID:       row.ID,
		Username: row.Username,
		Password: row.Password,
		Notification: domain.Notification{
			Email: row.Notification.Email,
			SMS:   row.Notification.SMS,
		},
	}
}

func (r userRepo) toRow(u *domain.User) userRow {
	return userRow{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Notification: notification{
			Email: u.Notification.Email,
			SMS:   u.Notification.SMS,
		},
	}
}
