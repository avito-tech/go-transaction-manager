module github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2

go 1.13

require (
	github.com/avito-tech/go-transaction-manager/v2 v2.0.0
	github.com/jackc/pgx/v5 v5.5.1
	github.com/pashagolub/pgxmock/v2 v2.12.0
	github.com/stretchr/testify v1.8.2
)

replace github.com/avito-tech/go-transaction-manager/v2 v2.0.0 => ../../