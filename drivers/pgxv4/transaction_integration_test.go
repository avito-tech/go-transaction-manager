//go:build with_real_db
// +build with_real_db

package pgxv4_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

func db(ctx context.Context) (*pgxpool.Pool, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"user", "pass", "localhost", 5432, "db",
	)

	pool, err := pgxpool.Connect(ctx, uri)
	if err != nil {
		return nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users_v4 (user_id SERIAL, username TEXT)`
	_, err = pool.Exec(ctx, sqlStmt)

	return pool, err
}

func TestTransaction_WithRealDB(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pool, err := db(ctx)
	require.NoError(t, err)
	defer pool.Close()

	f := pgxv4.NewDefaultFactory(pool)

	_, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	require.NoError(t, tr.Rollback(ctx))
	require.False(t, tr.IsActive())

	require.ErrorIs(t, tr.Commit(ctx), pgx.ErrTxClosed)
	require.ErrorIs(t, tr.Rollback(ctx), pgx.ErrTxClosed)
}
