package sqlx_test

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/transaction"
	trmcontext "github.com/avito-tech/go-transaction-manager/transaction/context"
	"github.com/avito-tech/go-transaction-manager/transaction/manager"
	"github.com/avito-tech/go-transaction-manager/transaction/settings"
)

// Example demonstrates a work of manager.ChainedMW.
func Example_chained() {
	// connect DB
	db1 := newDB()

	defer db1.Close() //nolint:errcheck

	db2 := newDB()

	defer db2.Close() //nolint:errcheck

	// create DB
	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err := db1.Exec(sqlStmt)
	checkErr(err, sqlStmt)

	_, err = db2.Exec(sqlStmt)
	checkErr(err, sqlStmt)

	// init manager
	ctxKey1 := trmcontext.Generate()
	m1 := manager.New(
		trmsqlx.NewDefaultFactory(db1),
		manager.WithSettings(settings.New(settings.WithCtxKey(ctxKey1))),
	)
	r1 := newRepo(db1, trmsqlx.NewCtxGetter(trmcontext.New(ctxKey1)))

	ctxKey2 := trmcontext.Generate()
	m2 := manager.New(
		trmsqlx.NewDefaultFactory(db2),
		manager.WithSettings(settings.New(settings.WithCtxKey(ctxKey2))),
	)
	r2 := newRepo(db2, trmsqlx.NewCtxGetter(trmcontext.New(ctxKey2)))

	chainedManager := manager.NewChained([]transaction.Manager{m1, m2})

	u := &user{
		Username: "username",
	}

	ctx := context.Background()

	err = chainedManager.Do(ctx, func(ctx context.Context) error {
		if err := r1.Save(ctx, u); err != nil {
			return err
		}

		if err := r2.Save(ctx, u); err != nil {
			return err
		}

		if err := chainedManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			if err = r1.Save(ctx, u); err != nil {
				return err
			}

			return r2.Save(ctx, u)
		}); err != nil {
			return err
		}

		return nil
	})
	checkErr(err)

	userFromDB1, err := r1.GetByID(ctx, u.ID)
	checkErr(err)

	userFromDB2, err := r1.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB1, userFromDB2)

	// Output: &{1 new_username} &{1 new_username}
}
