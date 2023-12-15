package with_or_without_trm

import (
	"context"

	"github.com/jmoiron/sqlx"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

type trmRepo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func newTrmRepo(db *sqlx.DB, c *trmsqlx.CtxGetter) *trmRepo {
	return &trmRepo{db: db, getter: c}
}

func (r *trmRepo) GetByID(ctx context.Context, id int64) (*user, error) {
	query := "SELECT * FROM user WHERE user_id = ?;"
	u := user{}

	return &u, r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &u, r.db.Rebind(query), id)
}

func (r *trmRepo) Save(ctx context.Context, u *user) error {
	query := `UPDATE user SET username = :username WHERE user_id = :user_id;`
	if u.ID == 0 {
		query = `INSERT INTO user (username) VALUES (:username);`
	}

	res, err := sqlx.NamedExecContext(ctx, r.getter.DefaultTrOrDB(ctx, r.db), r.db.Rebind(query), u)
	if err != nil {
		return err
	} else if u.ID != 0 {
		return nil
	} else if u.ID, err = res.LastInsertId(); err != nil {
		return err
	}

	return err
}
