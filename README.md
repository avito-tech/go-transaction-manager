# Go transaction manager

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
<!-- #TODO add images
[![GoDoc][doc-img]][doc] [![Coverage Status][cov-img]][cov] ![test][test-img])
-->

Transaction manager is an abstraction to coordinate database transaction boundaries.

## Supported implementations

* [sqlx](https://github.com/jmoiron/sqlx) (Go 1.16)

<!-- #TODO: 
* [sql](https://pkg.go.dev/database/sql) (Go 1.16)
* [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) (Go 1.16)
-->

## Installation

```bash
go get github.com/avito-tech/go-transaction-manager
```

### Backwards Compatibility

The library is compatible with the most recent two versions of Go.
Compatibility beyond that is not guaranteed.

## Usage

Below is an example how to start transaction. Check [example_test.go](sqlx/example_test.go) for more usage.

<!-- #TODO: add example -->
```go

```