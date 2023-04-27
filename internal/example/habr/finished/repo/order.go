package repo

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"

	"finish/domain"
)

type orderRow struct {
	ID        domain.OrderID   `db:"id"`
	ProductID domain.ProductID `db:"product_id"`
	UserID    domain.UserID    `db:"user_id"`
	Quantity  int64            `db:"quantity"`
}

type orderRepo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewOrderRepo(db *sqlx.DB, getter *trmsqlx.CtxGetter) *orderRepo {
	return &orderRepo{db: db, getter: getter}
}

func (r *orderRepo) GetByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	query := `SELECT * FROM "order" WHERE id = ?;`
	row := orderRow{}

	err := r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &row, r.db.Rebind(query), id)
	if err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *orderRepo) GetByUserID(ctx context.Context, id domain.UserID) (*domain.Order, error) {
	query := `SELECT * FROM "order" WHERE user_id = ?;`
	row := orderRow{}

	err := r.getter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &row, r.db.Rebind(query), id)
	if err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *orderRepo) Save(ctx context.Context, o *domain.Order) error {
	query := `INSERT INTO "order" (product_id, user_id, quantity)
VALUES (:product_id, :user_id, :quantity)
ON CONFLICT (id) DO UPDATE SET product_id = excluded.product_id,
                               quantity   = excluded.quantity
RETURNING id`

	rows, err := sqlx.NamedQueryContext(
		ctx,
		r.getter.DefaultTrOrDB(ctx, r.db),
		r.db.Rebind(query),
		r.toRow(o),
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}

	err = rows.Scan(&o.ID)

	return err
}

func (r orderRepo) toModel(row orderRow) *domain.Order {
	return &domain.Order{
		ID:        row.ID,
		ProductID: row.ProductID,
		UserID:    row.UserID,
		Quantity:  row.Quantity,
	}
}

func (r orderRepo) toRow(u *domain.Order) orderRow {
	return orderRow{
		ID:        u.ID,
		ProductID: u.ProductID,
		UserID:    u.UserID,
		Quantity:  u.Quantity,
	}
}
