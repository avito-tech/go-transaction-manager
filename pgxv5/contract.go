//go:build go1.16
// +build go1.16

// Package pgxv5 is an implementation of trm.Transaction interface by Transaction for pgx.Tx.
package pgxv5

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Tr is an interface to work with pgx.Conn, pgxpool.Conn or pgxpool.Pool
// StmtContext and Stmt are not implemented!
type Tr interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// Transactional is an interface to work with pgx.Conn, pgxpool.Conn or pgxpool.Pool.
type Transactional interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}
