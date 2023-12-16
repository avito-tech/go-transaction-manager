#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

echo "\ntrm"
cd trm && golangci-lint run -c $ROOT/.golangci.yml --timeout=2m $@
cd $ROOT

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && golangci-lint run -c $ROOT/.golangci.yml --timeout=2m $@

    cd $ROOT
  fi
done