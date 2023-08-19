# Go transaction manager

[![Go Reference](https://pkg.go.dev/badge/github.com/avito-tech/go-transaction-manager.svg)](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager)
[![Test Status](https://github.com/avito-tech/go-transaction-manager/actions/workflows/main.yaml/badge.svg)](https://github.com/avito-tech/go-transaction-manager/actions?query=branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/avito-tech/go-transaction-manager/badge.svg?branch=main)](https://coveralls.io/github/avito-tech/go-transaction-manager?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/avito-tech/go-transaction-manager)](https://goreportcard.com/report/github.com/avito-tech/go-transaction-manager)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Transaction manager is an abstraction to coordinate database transaction boundaries.

Easiest way to get the perfect repository.

## Supported implementations

* [database/sql](https://pkg.go.dev/database/sql), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/sql) (Go 1.13)
* [jmoiron/sqlx](https://github.com/jmoiron/sqlx), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/sqlx) (Go 1.13)
* [gorm](https://github.com/go-gorm/gorm), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/gorm) (Go 1.16)
* [mongo-go-driver](https://github.com/mongodb/mongo-go-driver), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/mongo) (Go 1.13)
* [go-redis/redis](https://github.com/go-redis/redis), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/redis) (Go 1.17)
* [pgx_v4](https://github.com/jackc/pgx), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/pgx) (Go 1.16)
* [pgx_v5](https://github.com/jackc/pgx), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/pgx) (Go 1.16)

## Installation

```bash
go get github.com/avito-tech/go-transaction-manager
```

### Backwards Compatibility

The library is compatible with the most recent two versions of Go.
Compatibility beyond that is not guaranteed.

## Usage

**To use multiple transactions from different databases**, you need to set CtxKey in [Settings](trm/settings.go) by [WithCtxKey](trm/settings/option.go).

**For nested transactions with different transaction managers**, you need to use [ChainedMW](trm/manager/chain.go) ([docs](https://pkg.go.dev/github.com/github.com/avito-tech/go-transaction-manager)).

**To skip a transaction rollback due to an error, use [ErrSkip](trm/manager.go#L20) or [Skippable](trm/manager.go#L24)**

### Explanation of the approach ([English](https://www.youtube.com/watch?v=aRsea6FFAyA), [Russian](https://habr.com/ru/companies/avito/articles/727168/))

### Examples with an ideal *repository* and nested transactions.

* [database/sql](sql/example_test.go)
* [jmoiron/sqlx](sqlx/example_test.go)
* [gorm](gorm/example_test.go)
* [mongo-go-driver](mongo/example_test.go)
* [go-redis/redis](redis/example_test.go)
* [pgx_v4](pgxv4/example_test.go)
* [pgx_v5](pgxv5/example_test.go)


Below is an example how to start usage.

```go
package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

func main() {
	db, err := sqlx.Open("sqlite3", "file:test?mode=memory")
	checkErr(err)
	defer db.Close()

	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err = db.Exec(sqlStmt)
	checkErr(err, sqlStmt)

	r := newRepo(db, trmsqlx.DefaultCtxGetter)
	ctx := context.Background()
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	u := &user{Username: "username"}

	err = trManager.Do(ctx, func(ctx context.Context) error {
		checkErr(r.Save(ctx, u))

		return trManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"
			return r.Save(ctx, u)
		})
	})
	checkErr(err)

	userFromDB, err := r.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB)
}

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}

type repo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func newRepo(db *sqlx.DB, c *trmsqlx.CtxGetter) *repo {
	return &repo{db: db, getter: c}
}

type user struct {
	ID       int64  `db:"user_id"`
	Username string `db:"username"`
}

func (r *repo) GetByID(ctx context.Context, id int64) (*user, error) {
	query := "SELECT * FROM user WHERE user_id = ?;"
	u := user{}

	return &u, r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &u, r.db.Rebind(query), id)
}

func (r *repo) Save(ctx context.Context, u *user) error {
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
```

## Benchmark

[Comparing](internal/benchmark/with_or_without_trm/README.md) examples with and without trm.