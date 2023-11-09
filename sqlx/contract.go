// Package sqlx is an implementation of trm.Transaction interface by Transaction for sqlx.Tx.
package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tr is an interface to work with sqlx.DB or sqlx.Tx.
// Stmtx, StmtxContext, NamedStmt and NamedStmtContext are not implemented!
//
//nolint:interfacebloat
type Tr interface {
	sqlx.ExtContext

	sqlx.Preparer
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	sqlx.Execer
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	sqlx.Queryer
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
