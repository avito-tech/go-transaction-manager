//go:build go1.19 && with_real_db
// +build go1.19,with_real_db

package pgxv5_test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/v2/manager"
)

// Example demonstrates the implementation of the Repository pattern by trm.Manager.
func Example() {
	ctx := context.Background()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"user", "pass", "localhost", 5432, "db",
	)

	pool, err := pgxpool.New(ctx, uri)
	checkErr(err)

	defer pool.Close()

	sqlStmt := `CREATE TABLE IF NOT EXISTS users_v5 (user_id SERIAL, username TEXT)`
	_, err = pool.Exec(ctx, sqlStmt)
	checkErr(err, sqlStmt)

	r := newRepo(pool, trmpgx.DefaultCtxGetter)
	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))

	u := &user{
		Username: "username",
	}

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
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func newRepo(db *pgxpool.Pool, c *trmpgx.CtxGetter) *repo {
	repo := &repo{
		db:     db,
		getter: c,
	}

	return repo
}

type user struct {
	ID       int64
	Username string
}

func (r *repo) GetByID(ctx context.Context, id int64) (*user, error) {
	query := `SELECT * FROM users_v5 WHERE user_id=$1`

	conn := r.getter.DefaultTrOrDB(ctx, r.db)
	row := conn.QueryRow(ctx, query, id)

	user := &user{}

	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) Save(ctx context.Context, u *user) error {
	isNew := u.ID == 0
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	if !isNew {
		query := `UPDATE users_v5 SET username = $1 WHERE user_id = $2`

		if _, err := conn.Exec(ctx, query, u.Username, u.ID); err != nil {
			return err
		}

		return nil
	}

	query := `INSERT INTO users_v5 (username) VALUES ($1) RETURNING user_id`

	err := conn.QueryRow(ctx, query, u.Username).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
