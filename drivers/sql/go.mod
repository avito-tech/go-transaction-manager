module github.com/avito-tech/go-transaction-manager/drivers/sql/v2

go 1.25.0

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.2
	github.com/golang/mock v1.6.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/stretchr/testify v1.11.1
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// go mod edit -replace=github.com/avito-tech/go-transaction-manager/trm/v2=../../
// go get github.com/avito-tech/go-transaction-manager/trm/v2
// go mod edit -dropreplace=github.com/avito-tech/go-transaction-manager/trm/v2

// replace github.com/avito-tech/go-transaction-manager/trm/v2 => ../../trm
