package with_or_without_trm

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func BenchmarkClean_MockDB(b *testing.B) {
	db, dbmock := dbMock(b)

	runClean(b, db)

	require.NoError(b, dbmock.ExpectationsWereMet())
}

func BenchmarkClean_SQLite_File(b *testing.B) {
	db := sqlite3(b, "file:test_clean.dbs")
	defer db.Close()

	runClean(b, db)
}

func BenchmarkClean_SQLite_Memory(b *testing.B) {
	db := sqlite3(b, "file:test?mode=memory")
	defer db.Close()

	runClean(b, db)
}

func runClean(b *testing.B, db *sqlx.DB) {
	r := newRepo(db)

	for i := 0; i < b.N; i++ {
		err := func() error {
			tx, err := db.Beginx()
			require.NoError(b, err)

			defer func() {
				if err != nil {
					tx.Rollback() //nolint:gosec
				}
			}()

			u := &user{Username: "username"}
			require.NoError(b, r.Save(tx, u))

			func() {
				hasTx := true
				if tx == nil {
					hasTx = true
					tx, err = db.Beginx()
					require.NoError(b, err)

					defer func() {
						if err != nil {
							tx.Rollback() //nolint:gosec
						}
					}()
				}

				u.Username = "new_username"
				require.NoError(b, r.Save(tx, u))

				if !hasTx {
					err = tx.Commit()
				}
			}()

			err = tx.Commit()

			return err
		}()
		require.NoError(b, err)
	}
}
