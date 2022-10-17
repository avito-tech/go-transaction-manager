package sqlx_test

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/transaction/manager"
)

func newDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "file:test?mode=memory")
	checkErr(err)

	return db
}

type repo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func newRepo(db *sqlx.DB, c *trmsqlx.CtxGetter) *repo {
	return &repo{
		db:     db,
		getter: c,
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

	err := r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &row, r.db.Rebind(query), id)
	if err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *repo) Save(ctx context.Context, u *user) error {
	isNew := u.ID == 0

	query := `UPDATE user SET username = :username WHERE user_id = :user_id;`
	if isNew {
		query = `INSERT INTO user (username) VALUES (:username);`
	}

	res, err := sqlx.NamedExecContext(
		ctx,
		r.getter.DefaultTrOrDB(ctx, r.db),
		r.db.Rebind(query),
		r.toRow(u),
	)
	if err != nil {
		return err
	} else if !isNew {
		return nil
	} else if u.ID, err = res.LastInsertId(); err != nil {
		return err
	}

	// For PostgreSql need to use NamedQueryContext with RETURNING
	// DO UPDATE SET username = EXCLUDED.username RETURNING id;
	// defer res.Next()
	// if u.ID == 0 && res.Next() {
	//		if err = res.Scan(&u.ID); err != nil {
	//			return err
	//		}
	//	}

	return nil
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

// Example demonstrates the implementation of the Repository pattern by TrManager.
func Example() {
	db := newDB()

	defer db.Close() //nolint:errcheck

	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err := db.Exec(sqlStmt)
	checkErr(err, sqlStmt)

	r := newRepo(db, trmsqlx.DefaultCtxGetter)

	u := &user{
		Username: "username",
	}

	ctx := context.Background()
	trManager := manager.New(trmsqlx.NewFactory(db))

	err = trManager.Do(ctx, func(ctx context.Context) error {
		if err := r.Save(ctx, u); err != nil {
			return err
		}

		if err := trManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			return r.Save(ctx, u)
		}); err != nil {
			return err
		}

		return nil
	})
	checkErr(err)

	userFromDB, err := r.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB)

	// Output: &{1 new_username}
}

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
