package uow_test

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
	trmuow "github.com/avito-tech/go-transaction-manager/uow"
)

func Example() {
	db := newDB()

	defer db.Close() //nolint:errcheck

	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err := db.Exec(sqlStmt)
	checkErr(err, sqlStmt)

	r := newRepo(db)

	u := &user{
		Username: "username",
	}

	ctx := context.Background()
	s := settings.Must(settings.WithCtxKey(trmuow.DefaultCtxKey))
	s2 := settings.Must(settings.WithCtxKey(trmuow.DefaultCtxKey))
	_ = s2

	uowTrManager := manager.Must(
		trmuow.NewDefaultFactory(
			manager.Must(trmsqlx.NewDefaultFactory(db)),
		),
		manager.WithSettings(
			s),
	)

	err = uowTrManager.Do(ctx, func(ctx context.Context) error {
		if err := r.Save(ctx, u); err != nil {
			return err
		}

		err = uowTrManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			return r.Save(ctx, u)
		})
		if err != nil {
			return err
		}

		_, err = r.GetByID(ctx, 1)
		fmt.Println("id:", u.ID, ", error:", err)

		return nil
	})
	checkErr(err)

	userFromDB, err := r.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB)

	// Output: id: 0 , error: sql: no rows in result set
	// &{1 new_username}
}

func newDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test?mode=memory")
	checkErr(err)

	return db
}

type repo struct {
	db *sqlx.DB
}

func newRepo(db *sqlx.DB) *repo {
	return &repo{
		db: db,
	}
}

type user struct {
	ID       int64
	Username string
}

type userRow struct {
	ID       int64  `db:"user_id"`
	Username string `db:"username"`
}

func (r *repo) GetByID(ctx context.Context, id int64) (*user, error) {
	query := "SELECT * FROM user WHERE user_id = ?;"

	row := userRow{}

	err := trmsqlx.DefaultCtxGetter.
		DefaultTrOrDB(ctx, r.db).
		GetContext(ctx, &row, r.db.Rebind(query), id)
	if err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *repo) Save(ctx context.Context, u *user) error {
	_, err := trmuow.DefaultCtxGetter.DefaultTr(ctx).Register(ctx, func(ctx context.Context) (interface{}, error) {
		isNew := u.ID == 0

		query := `UPDATE user SET username = :username WHERE user_id = :user_id;`
		if isNew {
			query = `INSERT INTO user (username) VALUES (:username);`
		}

		res, err := sqlx.NamedExecContext(
			ctx,
			trmsqlx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db),
			r.db.Rebind(query),
			r.toRow(u),
		)
		if err != nil {
			return nil, err
		} else if !isNew {
			return u, nil
		} else if u.ID, err = res.LastInsertId(); err != nil {
			return nil, err
		}

		// For PostgreSql need to use NamedQueryContext with RETURNING
		// DO UPDATE SET username = EXCLUDED.username RETURNING id;
		// defer res.Next()
		// if u.ID == 0 && res.Next() {
		//		if err = res.Scan(&u.ID); err != nil {
		//			return err
		//		}
		//	}

		return u, nil
	})

	return err
}

func (r *repo) toRow(model *user) userRow {
	return userRow{
		ID:       model.ID,
		Username: model.Username,
	}
}

func (r *repo) toModel(row userRow) *user {
	return &user{
		ID:       row.ID,
		Username: row.Username,
	}
}

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
