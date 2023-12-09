module github.com/avito-tech/go-transaction-manager/db/pgxv4/v2

go 1.13

require (
	github.com/avito-tech/go-transaction-manager/v2 v2.0.0
	github.com/jackc/pgconn v1.14.1
	github.com/jackc/pgx/v4 v4.18.1
	github.com/pashagolub/pgxmock v1.8.0
	github.com/stretchr/testify v1.8.2
)

replace github.com/avito-tech/go-transaction-manager/v2 v2.0.0 => ../../
