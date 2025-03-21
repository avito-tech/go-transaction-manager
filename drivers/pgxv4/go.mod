module github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2

go 1.16

require (
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0-rc10
	github.com/jackc/pgconn v1.14.2
	github.com/jackc/pgx/v4 v4.18.1
	github.com/lib/pq v1.10.9 // indirect
	github.com/pashagolub/pgxmock v1.8.0
	github.com/stretchr/testify v1.8.2
	go.uber.org/goleak v1.3.0
	golang.org/x/crypto v0.14.0 // indirect
)
