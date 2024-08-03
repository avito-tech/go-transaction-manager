//go:build with_real_db
// +build with_real_db

package pgxv5_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

func db(ctx context.Context) (*pgxpool.Pool, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"user", "pass", "localhost", 5432, "db",
	)

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users_v5 (user_id SERIAL, username TEXT)`
	_, err = pool.Exec(ctx, sqlStmt)

	return pool, err
}

func TestTransaction_WithRealDB(t *testing.T) {
	ctx := context.Background()

	pool, err := db(ctx)
	require.NoError(t, err)
	defer pool.Close()

	f := pgxv5.NewDefaultFactory(pool)

	_, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	require.NoError(t, tr.Rollback(ctx))
	require.False(t, tr.IsActive())

	require.ErrorIs(t, tr.Commit(ctx), pgx.ErrTxClosed)
	require.ErrorIs(t, tr.Rollback(ctx), pgx.ErrTxClosed)
}

// transaction should release all resources if context is cancelled
// otherwise pool.Close() is blocked forever
func TestTransaction_WithRealDB_RollbackOnContextCancel(t *testing.T) {
	ctx := context.Background()

	pool, err := db(ctx)
	require.NoError(t, err)

	defer func() {
		waitPoolIsClosed(t, pool)
	}()

	f := pgxv5.NewDefaultFactory(pool)

	ctx, cancel := context.WithCancel(ctx)

	_, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	require.True(t, tr.IsActive())

	cancel()
}

func waitPoolIsClosed(t *testing.T, pool *pgxpool.Pool) {
	const checkTick = 50 * time.Millisecond
	const waitDurationDeadline = 30 * time.Second

	var poolClosed atomic.Bool
	poolClosed.Store(false)

	go func() {
		pool.Close()
		poolClosed.Store(true)
	}()

	require.Eventually(
		t,
		func() bool {
			return poolClosed.Load()
		},
		waitDurationDeadline,
		checkTick)

	// https://github.com/jackc/pgx/issues/1641
	// pool triggerHealthCheck leaves stranded goroutines for 500ms
	// otherwise goleak error is triggered
	const waitPoolHealthCheck = 500 * time.Millisecond
	time.Sleep(waitPoolHealthCheck)
}
