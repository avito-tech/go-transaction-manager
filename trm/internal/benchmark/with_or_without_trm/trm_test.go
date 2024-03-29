package with_or_without_trm

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func BenchmarkTRM_MockDB(b *testing.B) {
	db, dbmock := dbMock(b)

	runTRM(b, db)

	require.NoError(b, dbmock.ExpectationsWereMet())
}

func BenchmarkTRM_SQLite_File(b *testing.B) {
	db := sqlite3(b, "file:test_trm.db")
	defer db.Close()

	runTRM(b, db)
}

func BenchmarkTRM_SQLite_Memory(b *testing.B) {
	db := sqlite3(b, "file:test?mode=memory")
	defer db.Close()

	runTRM(b, db)
}

func runTRM(b *testing.B, db *sqlx.DB) {
	r := newTrmRepo(db, trmsqlx.DefaultCtxGetter)
	ctx := context.Background()
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))

	for i := 0; i < b.N; i++ {
		u := &user{Username: "username"}

		err := trManager.Do(ctx, func(ctx context.Context) error {
			require.NoError(b, r.Save(ctx, u))

			return trManager.Do(ctx, func(ctx context.Context) error {
				u.Username = "new_username"
				return r.Save(ctx, u)
			})
		})
		require.NoError(b, err)
	}
}
