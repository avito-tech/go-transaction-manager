#!/bin/bash

cd ../

ROOT=$(pwd)

go mod tidy

drivers=$($DIR/utils/drivers.sh)

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && go mod tidy

    cd $ROOT
  fi
done
