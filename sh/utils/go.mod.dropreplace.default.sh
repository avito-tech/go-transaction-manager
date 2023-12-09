#!/bin/bash

go mod edit -dropreplace=github.com/avito-tech/go-transaction-manager/v2@v2=../../

# https://github.com/golang/go/issues/51932
# https://go.dev/doc/tutorial/workspaces
