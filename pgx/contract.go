// Package pgx is an implementation of trm.Transaction interface by Transaction for pgx.Tx.
package pgx

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Tr is an interface to work with pgx.DB or pgx.Tx.
// StmtContext and Stmt are not implemented!
type Tr interface {
	//Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)

	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)

	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
