#!/bin/bash

DIR=$(pwd)

drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

golist() {
  go list ./... | grep -v mock | grep -v internal/
}

gotest() {
  cd $driver

  go test -mod=readonly $(golist) "$@"

  cd $ROOT
}

cd trm && go test $(golist) $@ &
cd $ROOT

for driver in $drivers; do
  if [ -d "$driver" ]; then
    gotest $@ &
  fi
done

wait
