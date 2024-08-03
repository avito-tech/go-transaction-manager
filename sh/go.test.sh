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

  go test -mod=readonly -race $(golist) "$@"

  local exit_code=$?

  cd $ROOT

  (exit $exit_code);
}

cd trm && go test $(golist) $@ &
cd $ROOT

pids=()
for driver in $drivers; do
  if [ -d "$driver" ]; then
    gotest $@ &
    pids+=($!)
  fi
done

exit_code=0
for pid in ${pids[*]}; do
    wait $pid || exit_code=1
done

exit $exit_code
