#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

echo "\ntrm"
cd trm && golangci-lint formatters
cd $ROOT

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && golangci-lint formatters

    cd $ROOT
  fi
done
