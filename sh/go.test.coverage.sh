#!/bin/bash

DIR=$(pwd)
ROOT="$(pwd)/.."

COVERAGE_TEST="-cover -covermode=atomic -test.gocoverdir=$ROOT/coverage"

drivers=("$($DIR/utils/drivers.sh)")

mkdir -p $ROOT/coverage

./go.test.sh $COVERAGE_TEST "$@"

go tool covdata textfmt -i="$ROOT/coverage" -o "$ROOT/coverage.out"