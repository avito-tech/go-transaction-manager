package repo

import (
	"github.com/jmoiron/sqlx"

	"init/domain"
)

type orderRow struct {
	ID        domain.OrderID   `db:"id"`
	ProductID domain.ProductID `db:"product_id"`
	UserID    domain.UserID    `db:"user_id"`
	Quantity  int64            `db:"quantity"`
}

type ordrRepo struct {
	db *sqlx.DB
}

func NewOrdrRepo(db *sqlx.DB) *ordrRepo {
	return &ordrRepo{db: db}
}

func (r *ordrRepo) GetByID(tr domain.Tr, id domain.OrderID) (*domain.Order, error) {
	query := `SELECT * FROM "order" WHERE id = ?;`
	row := orderRow{}

	if tr != nil {
		tr = r.db
	}

	if err := tr.Get(&row, r.db.Rebind(query), id); err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *ordrRepo) GetByUserID(tr domain.Tr, id domain.UserID) (*domain.Order, error) {
	query := `SELECT * FROM "order" WHERE user_id = ?;`
	row := orderRow{}

	if tr != nil {
		tr = r.db
	}

	if err := tr.Get(&row, r.db.Rebind(query), id); err != nil {
		return nil, err
	}

	return r.toModel(row), nil
}

func (r *ordrRepo) Save(tr domain.Tr, o *domain.Order) error {
	query := `INSERT INTO "order" (product_id, user_id, quantity)
VALUES (:product_id, :user_id, :quantity)
ON CONFLICT (id) DO UPDATE SET product_id = excluded.product_id,
                               quantity   = excluded.quantity
RETURNING id`

	if tr != nil {
		tr = r.db
	}

	rows, err := sqlx.NamedQuery(tr, query, r.toRow(o))
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

func (r ordrRepo) toModel(row orderRow) *domain.Order {
	return &domain.Order{
		ID:        row.ID,
		ProductID: row.ProductID,
		UserID:    row.UserID,
		Quantity:  row.Quantity,
	}
}

func (r ordrRepo) toRow(u *domain.Order) orderRow {
	return orderRow{
		ID:        u.ID,
		ProductID: u.ProductID,
		UserID:    u.UserID,
		Quantity:  u.Quantity,
	}
}
