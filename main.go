package main

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Data struct {
	Id   int
	Name string
}

type Repo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func main() {
	log.Printf("Demo")

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"test_user", "test_pass", "localhost", 5432, "test_db")

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, uri)
	if err != nil {
		log.Fatalf("Error psql connection: %v", err)
	}

	repo := newRepo(pool, trmpgx.DefaultCtxGetter)

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))

	err = trManager.Do(ctx, func(ctx context.Context) error {
		log.Printf("Called in the first TX")
		_, _ = repo.Read(ctx)
		return nil
	})

	data, err := repo.Read(ctx)
	if err != nil {
		log.Fatalf("err read from repo: %v", err)
	}

	log.Printf("data: %v", data)
}

func newRepo(db *pgxpool.Pool, c *trmpgx.CtxGetter) *Repo {
	repo := &Repo{
		db:     db,
		getter: c,
	}
	return repo
}

func (r *Repo) Read(ctx context.Context) ([]Data, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	rows, err := conn.Query(ctx, `SELECT * FROM my_table`)
	if err != nil {
		log.Fatalf("err query: %v", err)
	}
	defer rows.Close()

	res := []Data{}

	for rows.Next() {
		d := Data{}
		if err := rows.Scan(&d.Id, &d.Name); err != nil {
			return nil, err
		}
		res = append(res, d)
	}

	return res, nil
}

func (r *Repo) Save(ctx context.Context, d Data) error {
	return nil
}
