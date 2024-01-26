package test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewDBMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, dbmock, _ := sqlmock.New()
	// need to solve goroutine leak detection https://kumakichi.github.io/goroutine-leak.html
	t.Cleanup(func() {
		dbmock.ExpectClose()

		_ = db.Close()
	})
	return db, dbmock
}
