#!/bin/bash

cd ../

ROOT=$(pwd)

go mod vendor

drivers=$($DIR/utils/drivers.sh)

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && go mod vendor

    cd $ROOT
  fi
done
