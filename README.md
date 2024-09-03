# Go transaction manager

[![Go Reference](https://pkg.go.dev/badge/github.com/avito-tech/go-transaction-manager.svg)](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/trm/v2)
[![Test Status](https://github.com/avito-tech/go-transaction-manager/actions/workflows/main.yaml/badge.svg)](https://github.com/avito-tech/go-transaction-manager/actions?query=branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/avito-tech/go-transaction-manager/badge.svg?branch=main)](https://coveralls.io/github/avito-tech/go-transaction-manager?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/avito-tech/go-transaction-manager)](https://goreportcard.com/report/github.com/avito-tech/go-transaction-manager/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Transaction manager is an abstraction to coordinate database transaction boundaries.

Easiest way to get the perfect repository.

## Supported implementations

* [database/sql](https://pkg.go.dev/database/sql), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/sql/v2) (
  Go 1.13)
* [jmoiron/sqlx](https://github.com/jmoiron/sqlx), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2) (
  Go 1.13)
* [gorm](https://github.com/go-gorm/gorm), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/gorm/v2) (
  Go 1.18)
* [mongo-go-driver](https://github.com/mongodb/mongo-go-driver), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/mongo/v2) (
  Go 1.13)
* [go-redis/redis](https://github.com/go-redis/redis), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/goredis8/v2) (
  Go 1.17)
* [pgx_v4](https://github.com/jackc/pgx/tree/v4), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2) (
  Go 1.16)
* [pgx_v5](https://github.com/jackc/pgx), [docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2) (
  Go 1.19)

## Installation

```bash
go get github.com/avito-tech/go-transaction-manager/trm/v2
```

To install some support database use `go get github.com/avito-tech/go-transaction-manager/drivers/{name}`.

For example `go get github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2`.

### Backwards Compatibility

The library is compatible with the most recent two versions of Go.
Compatibility beyond that is not guaranteed.

The critical bugs are firstly solved for the most recent two Golang versions and then for older ones if it is simple.

#### Disclaimer: Keep your dependencies up to date, even indirect ones.

`go get -u && go mod tidy` helps you.

**Note**: The go-transaction-manager uses some old dependencies to support backwards compatibility for old versions of Go.

## Usage

**To use multiple transactions from different databases**, you need to set CtxKey in [Settings](trm/settings.go)
by [WithCtxKey](trm/settings/option.go) ([docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/trm/v2)).

**For nested transactions with different transaction managers**, you need to use [ChainedMW](trm/manager/chain.go) ([docs](https://pkg.go.dev/github.com/avito-tech/go-transaction-manager/trm/v2/manager)).

**To skip a transaction rollback due to an error, use [ErrSkip](manager.go#L20) or [Skippable](manager.go#L24)**

### Explanation of the approach [English](https://www.youtube.com/watch?v=aRsea6FFAyA), Russian [article](https://habr.com/ru/companies/avito/articles/727168/) and [youtube](https://www.youtube.com/watch?v=fcdckM5sUxA).

### Examples with an ideal *repository* and nested transactions.

* [database/sql](drivers/sql/example_test.go)
* [jmoiron/sqlx](drivers/sqlx/example_test.go)
* [gorm](drivers/gorm/example_test.go)
* [mongo-go-driver](drivers/mongo/example_test.go)
* [go-redis/redis](drivers/goredis8/example_test.go)
* [pgx_v4](drivers/pgxv4/example_test.go)
* [pgx_v5](drivers/pgxv5/example_test.go)

Below is an example how to start usage.

```go
package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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

		// example of nested transactions
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

[Comparing](trm/internal/benchmark/with_or_without_trm/README.md) examples with and without trm.

## Contribution

### Requirements

- [golangci-lint](https://golangci-lint.run/welcome/install/)
- [make](https://www.gnu.org/software/make/#download)

### Local Running

* To install all dependencies use `make go.mod.tidy` or `make go.mod.vendor`.
* To run all tests use `make go.test` or `make go.test.with_real_db` for integration tests.

To run database by docker, there is [docker-compose.yaml](trm/drivers/test/docker-compose.yaml).
```bash
docker compose -f trm/drivers/test/docker-compose.yaml up
```

For full GitHub Actions run, you can use [act](https://github.com/nektos/act).

#### Running old go versions 

To stop Golang upgrading set environment variable `GOTOOLCHAIN=local` .

```sh
go install go1.16 # or older version
go1.16 install
```

Use `-mod=readonly` to prevent `go.mod` modification.

To run tests
```
go1.16 test -race -mod=readonly ./...
```

### How to bump up Golang version in CI/CD

1. Changes in [.github/workflows/main.yaml](.github/workflows/main.yaml).
   1. Add all old version of Go in `go-version:` for `tests-units` job.
   2. Update `go-version:` on current version of Go for `lint` and `tests-integration` jobs.
2. Update build tags by replacing `build go1.xx` on new version.


### Resolve problems with old version of dependencies

To build `go.mod` compatible for old version use `go mod tidy -compat=1.13` ([docs](https://go.dev/ref/mod#go-mod-tidy)).

However, `--compat` doesn't always work correct and we need to set some library versions manually.

1. `go get go.uber.org/multierr@v1.9.0` in [trm](trm), [sql](drivers/sql), [sqlx](drivers/sqlx).
2. `go get github.com/mattn/go-sqlite3@v1.14.14` in [trm](trm), [sql](drivers/sql), [sqlx](drivers/sqlx).
3. `go get github.com/stretchr/testify@v1.8.2` in [trm](trm), [sql](drivers/sql), [sqlx](drivers/sqlx), [goredis8](drivers/goredis8), [mongo](drivers/mongo).
4. `go get github.com/jackc/pgconn@v1.14.2` in [pgxv4](drivers/pgxv4). Golang version was bumped up from 1.12 to 1.17 in pgconn v1.14.3.
5. `go get golang.org/x/text@v0.13.0` in [pgxv4](drivers/pgxv4).