# Go transaction manager

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
<!-- #TODO add images
[![GoDoc][doc-img]][doc] [![Coverage Status][cov-img]][cov] ![test][test-img])
-->

Transaction manager is an abstraction to coordinate database transaction boundaries.

## Supported implementations

* [sqlx](https://github.com/jmoiron/sqlx) (Go 1.10)

<!-- #TODO: 
* [sql](https://pkg.go.dev/database/sql) (Go 1.10)
* [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) (Go 1.10)
-->

## Installation

```bash
go get github.com/avito-tech/go-transaction-manager
```

### Backwards Compatibility

The library is compatible with the most recent two versions of Go.
Compatibility beyond that is not guaranteed.

## Usage

Below is an example how to start transaction. Check [example_test.go](sqlx/example_test.go) for more usage.


```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	trsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type user struct {
	ID       int64  `db:"user_id"`
	Username string `db:"username"`
}

func (r *repo) Save(ctx context.Context, u *user) error {
	isNew := u.ID == 0

	query := `UPDATE user SET username = :username WHERE user_id = :user_id;`
	if isNew {
		query = `INSERT INTO user (username) VALUES (:username);`
	}

	res, err := sqlx.NamedExecContext(ctx, trsqlx.TrFromCtx(ctx, r.db), r.db.Rebind(query), u)
	if err != nil {
		return err
	} else if !isNew {
		return nil
	} else if u.ID, err = res.LastInsertId(); err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := sqlx.Open("sqlite3", "file:test?mode=memory")
	checkErr(err)
	defer db.Close()

	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	r := &repo{db: db}

	ctx := context.Background()
	trManager := trsqlx.NewTransactionManager(db)
	u := &user{Username: "username"}

	err = trManager.Do(ctx, func(ctx context.Context) error {

		checkErr(r.Save(ctx, u))

		return trManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			return r.Save(ctx, u)
		})
	})
	checkErr(err)
	fmt.Println(u)
}

func checkErr(err error, args ...interface{}) {
	if err != nil {
		log.Fatal(append([]interface{}{err}, args...)...)
		return
	}
}
```