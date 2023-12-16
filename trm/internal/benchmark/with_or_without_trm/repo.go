package with_or_without_trm

import (
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func newRepo(db *sqlx.DB) *repo {
	return &repo{db: db}
}

func (r *repo) GetByID(tx *sqlx.Tx, id int64) (*user, error) {
	query := "SELECT * FROM user WHERE user_id = ?;"
	u := user{}

	return &u, r.txOrDB(tx).Get(&u, r.db.Rebind(query), id)
}

func (r *repo) Save(tx *sqlx.Tx, u *user) error {
	query := `UPDATE user SET username = :username WHERE user_id = :user_id;`
	if u.ID == 0 {
		query = `INSERT INTO user (username) VALUES (:username);`
	}

	res, err := sqlx.NamedExec(r.txOrDB(tx), r.db.Rebind(query), u)
	if err != nil {
		return err
	} else if u.ID != 0 {
		return nil
	} else if u.ID, err = res.LastInsertId(); err != nil {
		return err
	}

	return err
}

//nolint:ireturn
func (r *repo) txOrDB(tx *sqlx.Tx) tr {
	if tx != nil {
		return tx
	}

	return r.db
}

type tr interface {
	sqlx.Ext

	Get(dest interface{}, query string, args ...interface{}) error
}
