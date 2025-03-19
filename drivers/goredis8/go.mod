module github.com/avito-tech/go-transaction-manager/drivers/goredis8/v2

go 1.14

require (
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0-rc10
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-redis/redismock/v8 v8.11.5
	github.com/stretchr/testify v1.9.0
	go.uber.org/goleak v1.3.0
	golang.org/x/net v0.33.0 // indirect
)
