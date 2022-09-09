// Package sqlx is an implementation of transaction.Transaction interface by Transaction for sqlx.Tx.
package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tr is an interface to work with sqlx.DB or sqlx.Tx.
type Tr interface {
	sqlx.ExtContext

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}
