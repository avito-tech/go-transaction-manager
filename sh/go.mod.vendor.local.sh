#!/bin/bash

cd ../

ROOT=$(pwd)

go mod vendor

drivers=$($DIR/utils/drivers.sh)

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"

    cd $driver && \
    $ROOT/sh/utils/go.mod.replace.default.sh
    go mod vendor && \
    $ROOT/sh/utils/go.mod.dropreplace.default.sh

    cd $ROOT
  fi
done
