module github.com/avito-tech/go-transaction-manager/drivers/sql/v2

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0-rc9.2
	github.com/golang/mock v1.6.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/stretchr/testify v1.9.0
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.9.0
)

// go mod edit -replace=github.com/avito-tech/go-transaction-manager/trm/v2=../../
// go get github.com/avito-tech/go-transaction-manager/trm/v2
// go mod edit -dropreplace=github.com/avito-tech/go-transaction-manager/trm/v2

// replace github.com/avito-tech/go-transaction-manager/trm/v2 => ../../trm
