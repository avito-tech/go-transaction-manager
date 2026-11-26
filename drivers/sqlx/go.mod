module github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2

go 1.13

require (
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/avito-tech/go-transaction-manager/drivers/sql/v2 v2.0.1
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.5
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/stretchr/testify v1.9.0
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.9.0
)
