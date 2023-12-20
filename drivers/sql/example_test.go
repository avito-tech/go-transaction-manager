package sql_test

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"

	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

// Example demonstrates the implementation of the Repository pattern by trm.Manager.
func Example() {
	db, err := sql.Open("sqlite3", "file:test?mode=memory")
	checkErr(err)

	defer db.Close() //nolint:errcheck

	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err = db.Exec(sqlStmt)
	checkErr(err, sqlStmt)

	r := newRepo(db, trmsql.DefaultCtxGetter)

	u := &user{
		Username: "username",
	}

	ctx := context.Background()
	trManager := manager.Must(
		trmsql.NewDefaultFactory(db),
		manager.WithCtxManager(trmcontext.DefaultManager),
	)

	err = trManager.Do(ctx, func(ctx context.Context) error {
		if err := r.Save(ctx, u); err != nil {
			return err
		}

		return trManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			return r.Save(ctx, u)
		})
	})
	checkErr(err)

	userFromDB, err := r.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB)

	// Output: &{1 new_username}
}

type repo struct {
	db     *sql.DB
	getter *trmsql.CtxGetter
}

func newRepo(db *sql.DB, c *trmsql.CtxGetter) *repo {
	return &repo{
		db:     db,
		getter: c,
	}
}

type user struct {
	ID       int64
	Username string
}

func (r *repo) GetByID(ctx context.Context, id int64) (*user, error) {
	query := "SELECT * FROM user WHERE user_id = ?;"

	u := &user{}

	err := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *repo) Save(ctx context.Context, u *user) error {
	isNew := u.ID == 0

	args := []interface{}{
		sql.Named("username", u.Username),
	}
	query := `INSERT INTO user (username) VALUES (:username);`

	if !isNew {
		query = `UPDATE user SET username = :username WHERE user_id = :user_id;`

		args = append(args, sql.Named("user_id", u.ID))
	}

	res, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query, args...)
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

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
