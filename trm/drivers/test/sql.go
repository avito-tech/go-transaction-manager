// Package test provides a set of utilities for testing.
// You should not use this package in your application code.
// TODO try to use Deprecated for the package
package test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// NewDBMock returns a new sql.DB and sqlmock.
//
//nolint:ireturn
func NewDBMock() (*sql.DB, sqlmock.Sqlmock) {
	db, dbmock, _ := sqlmock.New()

	return db, dbmock
}

// NewDBMockWithClose returns a new sql.DB and sqlmock, and close it after test.
//
//nolint:ireturn
func NewDBMockWithClose(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, dbmock := NewDBMock()

	// need to solve goroutine leak detection https://kumakichi.github.io/goroutine-leak.html
	Cleanup(t, func() {
		dbmock.ExpectClose()

		require.NoError(t, db.Close())
	})

	return db, dbmock
}
