package repo

import (
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func SetDB(db *sqlx.DB) {
	DB = db
}

func WithTransaction(tx *sqlx.Tx, fn func(*sqlx.Tx) error) (err error) {
	hasExternalTx := true
	if tx == nil {
		tx, err = DB.Beginx()
		if err != nil {
			return err
		}
		hasExternalTx = false
	}

	defer func() {
		if !hasExternalTx && err != nil {
			tx.Rollback()
		}
	}()

	err = fn(tx) // call a usecase
	if err != nil {
		return err
	}

	if !hasExternalTx {
		err = tx.Commit() // or tr.Rollback()
	}

	return nil
}
