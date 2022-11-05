// Package sqlx is an implementation of trm.Transaction interface by Transaction for sqlx.Tx.
package sqlx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tr is an interface to work with sqlx.DB or sqlx.Tx.
// Stmtx, StmtxContext, NamedStmt and NamedStmtContext are not implemented!
type Tr interface {
	sqlx.ExtContext

	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)

	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)

	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
