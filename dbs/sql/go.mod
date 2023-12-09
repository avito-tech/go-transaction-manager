module github.com/avito-tech/go-transaction-manager/db/sql/v2

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/avito-tech/go-transaction-manager/v2 v2.0.0
	github.com/golang/mock v1.6.0
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/stretchr/testify v1.8.2
	go.uber.org/multierr v1.9.0
)

// go mod edit -replace=github.com/avito-tech/go-transaction-manager/v2=../../
// go get github.com/avito-tech/go-transaction-manager/v2
// go mod edit -dropreplace=github.com/avito-tech/go-transaction-manager/v2

// replace github.com/avito-tech/go-transaction-manager/v2 => ../../
