package with_or_without_trm

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

//nolint:ireturn
func dbMock(b *testing.B) (*sqlx.DB, sqlmock.Sqlmock) {
	db, dbmock, err := sqlmock.New()
	require.NoError(b, err)

	dbmock.MatchExpectationsInOrder(false)

	for i := 0; i < b.N; i++ {
		dbmock.ExpectBegin()

		dbmock.ExpectExec(".+").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 0))

		dbmock.ExpectExec(".+").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 0))

		dbmock.ExpectCommit()
	}

	return sqlx.NewDb(db, "sqlmock"), dbmock
}

func sqlite3(b *testing.B, dsn string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", dsn)
	require.NoError(b, err)

	sqlStmt := `CREATE TABLE IF NOT EXISTS user (user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT);`
	_, err = db.Exec(sqlStmt)
	require.NoError(b, err)

	return db
}
