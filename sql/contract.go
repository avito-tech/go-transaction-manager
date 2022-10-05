// Package sql is an implementation of transaction.Transaction interface by Transaction for sql.Tx.
package sql

import (
	"context"
	"database/sql"
)

// Tr is an interface to work with sql.DB or sql.Tx.
// StmtContext and Stmt are not implemented!
type Tr interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)

	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
}
