#!/bin/bash

DIR=$(pwd)

drivers=$($DIR/utils/drivers.sh)

echo $drivers

cd ../

ROOT=$(pwd)

golist() {
  go list ./... | grep -v mock | grep -v internal/
}

gotest() {
  cd $driver

  go test $(golist) "$@"

  cd $ROOT
}

go test $(golist) $@ &

for driver in $drivers; do
  if [ -d "$driver" ]; then
    gotest $@ &
  fi
done

wait
