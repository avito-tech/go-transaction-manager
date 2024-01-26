package test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewDBMock() (*sql.DB, sqlmock.Sqlmock) {
	db, dbmock, _ := sqlmock.New()

	return db, dbmock
}

func NewDBMockWithClose(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, dbmock := NewDBMock()

	// need to solve goroutine leak detection https://kumakichi.github.io/goroutine-leak.html
	t.Cleanup(func() {
		dbmock.ExpectClose()

		_ = db.Close()
	})

	return db, dbmock
}
