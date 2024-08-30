# Proving that TRM doesn't affect updates to newer version. 

1. Install pgxv5 driver with old pgx version `github.com/jackc/pgx/v5@v5.5.1`

   `GOWORK=off go get github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2@v2.0.0-rc9.2`
2. In `go.mod` we see `github.com/jackc/pgx/v5 v5.5.1 // indirect`.
3. Call `GOWORK=off go mod tidy && GOWORK=off go mod vendor` to install the old versions.
4. Then, we can update `pgx` manually and see in `go.mod` the last version `github.com/jackc/pgx/v5 v5.6.0 // indirect`. 
   
   `go get github.com/jackc/pgx/v5`
   or
   `go mod tidy`
5. `go test ./...` to run [example_test.go](example_test.go)