// Package sql is an implementation of transaction.Transaction interface by Transaction for sql.Tx.
package sql

import (
	"context"
	"database/sql"
)

// Tr  is an interface to work with sql.DB or sql.Tx.
// TODO add methods from sql.Tx.
type Tr interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
