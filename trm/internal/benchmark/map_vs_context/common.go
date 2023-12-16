package map_vs_context

import (
	"database/sql"
	"log"

	"github.com/jinzhu/copier"

	"github.com/avito-tech/go-transaction-manager/trm/v2/internal/benchmark/common"
)

func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", "file:test?mode=memory")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type creator func() *sql.Tx

func creatorEmpty() creator {
	tr := &sql.Tx{}

	return func() *sql.Tx {
		return tr
	}
}

func creatorCopy(db *sql.DB) creator {
	srcTr, err := db.Begin()
	common.CheckErr(err)

	return func() *sql.Tx {
		tr := &sql.Tx{}

		err := copier.CopyWithOption(tr, &srcTr, copier.Option{DeepCopy: true})
		common.CheckErr(err)

		return tr
	}
}

func creatorRealTransaction(db *sql.DB) creator {
	return func() *sql.Tx {
		tr, err := db.Begin()
		common.CheckErr(err)

		return tr
	}
}
