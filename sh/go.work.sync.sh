#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && go work sync &

    cd $ROOT
  fi
done

wait