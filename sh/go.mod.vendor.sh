#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

echo "\ntrm"
cd trm && go mod vendor
cd $ROOT

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && go mod vendor

    cd $ROOT
  fi
done
