//go:build with_real_db
// +build with_real_db

package pgxv4_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

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
	require.NoError(t, tr.Rollback(ctx)) // idempotent: returns nil when already closed
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

	f := pgxv4.NewDefaultFactory(pool)

	ctx, cancel := context.WithCancel(ctx)

	_, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	require.True(t, tr.IsActive())

	cancel()
	require.NoError(t, tr.Rollback(ctx))
}

// TestTransaction_WithRealDB_NoConcurrentAccessOnContextCancel verifies that
// cancelling a context while a query is in-flight does not cause concurrent
// pgx.Tx access (jackc/pgx#2332: "slow write timer already active" panic).
//
// On main (before fix): awaitDone goroutine calls tx.Rollback concurrently
// with the in-flight query, causing "conn busy" errors or panic.
// After fix: awaitDone is removed, no concurrent access occurs.
func TestTransaction_WithRealDB_NoConcurrentAccessOnContextCancel(t *testing.T) {
	ctx := context.Background()

	pool, err := db(ctx)
	require.NoError(t, err)
	defer waitPoolIsClosed(t, pool)

	f := pgxv4.NewDefaultFactory(pool)
	ctx, cancel := context.WithCancel(ctx)

	txCtx, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	pgxTx := tr.Transaction().(pgx.Tx)

	queryCh := make(chan error, 1)
	go func() {
		_, err := pgxTx.Exec(txCtx, "SELECT pg_sleep(1)")
		queryCh <- err
	}()

	// wait for query to reach the server
	time.Sleep(50 * time.Millisecond)

	// cancel while query is in-flight:
	// before fix: awaitDone fires and calls Rollback concurrently → conn busy / panic
	// after fix:  no awaitDone, no concurrent access
	cancel()

	queryErr := <-queryCh
	require.Error(t, queryErr, "query must fail due to context cancellation")

	// Connection may be closed by pgx after context cancellation.
	// The key guarantee: no panic from concurrent pgx.Tx access (jackc/pgx#2332).
	_ = tr.Rollback(context.Background())
	require.False(t, tr.IsActive())
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
