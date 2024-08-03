#!/bin/bash

DIR=$(pwd)

drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

golist() {
  go list ./... | grep -v mock | grep -v internal/
}

verlte() {
    printf '%s\n%s' "$1" "$2" | sort -C -V
}

gotest() {
  cd $driver

  local go_mod_ver=$(sed -En 's/^go (.*)$/\1/p' go.mod)
  local go_ver=$(go version | sed -n 's/.*go\([0-9.]*\).*/\1/p')

  local exit_code=0
  if verlte $go_mod_ver $go_ver; then

    go test -race -mod=readonly $(golist) "$@"
    exit_code=$?
  fi

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
