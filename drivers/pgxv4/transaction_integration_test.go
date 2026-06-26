//go:build with_real_db
// +build with_real_db

package pgxv4_test

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
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

	f := pgxv4.NewDefaultFactory(pool)

	ctx, cancel := context.WithCancel(ctx)

	_, tr, err := f(ctx, settings.Must())
	require.NoError(t, err)

	require.True(t, tr.IsActive())

	cancel()
	require.ErrorIs(t, tr.Rollback(ctx), context.Canceled)
	require.ErrorIs(t, tr.Rollback(ctx), pgx.ErrTxClosed)
}

// TestTransaction_WithRealDB_NoDataRaceOnContextCancelDuringQuery_139 reproduces
// https://github.com/avito-tech/go-transaction-manager/issues/139.
//
// before fix: awaitDone goroutine calls tx.Rollback concurrently
// with the in-flight query, causing panic `panic: BUG: slow write timer already active`.
// https://github.com/jackc/pgx/blob/v5.10.0/pgconn/pgconn.go#L2115.
//
// pgx.Tx is not safe for concurrent use (jackc/pgx#2332). It causes panic when we call two commands simultaneously.
//
// cancelAfter controls that the transaction is canceled exactly when we run SQL query.
// cancelAfter should be less than pg_sleep_for.
// pg_sleep_for controls that the query is still running while the transaction is being canceled by cancelAfter.
func TestTransaction_WithRealDB_NoDataRaceOnContextCancelDuringQuery_139(t *testing.T) {
	ctx := context.Background()

	pool, err := db(ctx)
	require.NoError(t, err)

	defer waitPoolIsClosed(t, pool)

	trManager := manager.Must(pgxv4.NewDefaultFactory(pool))

	// 8 MB parameter writing keeps the connection write-busy to hit every phase of the protocol.
	// That forces the slow write timer to already be active, triggering a panic
	// even if we don't run without the race detector (-race).
	payload := strings.Repeat("x", 8*1024*1024)

	const (
		attempts          = 25
		explanationErrMsg = "Change cancelAfter or pg_sleep_for."
	)

	for i := 0; i < attempts; i++ {
		cancelAfter := time.Duration(1+2*i) * time.Millisecond
		ctx, cancel := context.WithCancel(ctx)

		err := trManager.Do(ctx, func(ctx context.Context) error {
			go func() {
				// cancel context when pgx executes a query.
				time.Sleep(cancelAfter)
				cancel()
			}()

			require.NoError(t, ctx.Err(), explanationErrMsg)

			_, err := pgxv4.DefaultCtxGetter.DefaultTrOrDB(ctx, pool).
				Exec(ctx, "SELECT pg_sleep_for('0.1 seconds'), length($1)", payload)

			require.ErrorIs(t, ctx.Err(), context.Canceled, explanationErrMsg)

			return err
		})

		require.ErrorIs(t, err, context.Canceled, explanationErrMsg)
	}
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
