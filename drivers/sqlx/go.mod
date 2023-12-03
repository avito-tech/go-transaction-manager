module github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2

go 1.13

require (
	github.com/avito-tech/go-transaction-manager/v2  v2.0.0
	github.com/avito-tech/go-transaction-manager/drivers/sql/v2  v2.0.0
)

replace github.com/avito-tech/go-transaction-manager/v2 v2.0.0 => ../../
replace github.com/avito-tech/go-transaction-manager/drivers/sql/v2 v2.0.0 => ../sql