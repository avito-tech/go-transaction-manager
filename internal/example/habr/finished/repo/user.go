package repo

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"

	"finish/domain"
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
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewUserRepo(db *sqlx.DB, getter *trmsqlx.CtxGetter) *userRepo {
	return &userRepo{
		db:     db,
		getter: getter,
	}
}
func (r *userRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query := `SELECT * FROM "user" WHERE id = ?;`
	row := userRow{}

	err := r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &row, r.db.Rebind(query), id)
	if err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *userRepo) Save(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO "user" (username, password, notification)
VALUES (:username, :password, :notification)
ON CONFLICT (id)
    DO UPDATE SET username     = excluded.username,
                  password     = excluded.password,
                  notification = excluded.notification
RETURNING id;`

	rows, err := sqlx.NamedQueryContext(
		ctx,
		r.getter.DefaultTrOrDB(ctx, r.db),
		r.db.Rebind(query),
		r.toRow(u),
	)
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
