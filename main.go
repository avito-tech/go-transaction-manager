package main

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	settings2 "github.com/avito-tech/go-transaction-manager/trm/settings"
	"github.com/jackc/pgx/v4"
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

	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		log.Fatalf("Error psql connection conn: %v", err)
	}
	_ = conn.Ping(ctx)

	pool, err := pgxpool.Connect(ctx, uri)
	if err != nil {
		log.Fatalf("Error psql connection pool: %v", err)
	}

	repo := newRepo(pool, trmpgx.DefaultCtxGetter)

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))
	//trManager := manager.Must(trmpgx.NewDefaultFactory(conn))

	settings := trmpgx.MustSettings(
		settings2.Must(),
		trmpgx.WithTxOptions(&pgx.TxOptions{IsoLevel: pgx.ReadCommitted}),
	)

	d := Data{Id: 1, Name: "111111111111"}
	if err := repo.Save(ctx, d); err != nil {
		log.Fatalf("err save: %v", err)
	}

	err = trManager.DoWithSettings(ctx, settings, func(ctx context.Context) error {

		d := Data{Id: 2, Name: "222222222222"}
		if err := repo.Save(ctx, d); err != nil {
			return err
		}

		log.Printf("Called in the first TX")
		//return errors.New("111")
		return nil
	})

	d = Data{Id: 3, Name: "333333333333"}
	if err := repo.Save(ctx, d); err != nil {
		log.Fatalf("err save: %v", err)
	}

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
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	query := `INSERT INTO my_table (id, name) VALUES ($1, $2)`
	if _, err := conn.Exec(ctx, query, d.Id, d.Name); err != nil {
		return err
	}

	return nil
}
