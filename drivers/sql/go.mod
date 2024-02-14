module github.com/avito-tech/go-transaction-manager/drivers/sql/v2

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.5.1
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0-rc6
	github.com/golang/mock v1.6.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/stretchr/testify v1.8.2
	go.uber.org/goleak v1.3.0 // indirect
	go.uber.org/multierr v1.9.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// go mod edit -replace=github.com/avito-tech/go-transaction-manager/trm/v2=../../
// go get github.com/avito-tech/go-transaction-manager/trm/v2
// go mod edit -dropreplace=github.com/avito-tech/go-transaction-manager/trm/v2

// replace github.com/avito-tech/go-transaction-manager/trm/v2 => ../../trm
