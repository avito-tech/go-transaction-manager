module github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.5.1
	github.com/avito-tech/go-transaction-manager/drivers/sql/v2 v2.0.0-rc9.1
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0-rc10
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.14
	github.com/stretchr/testify v1.8.2
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.9.0
)
