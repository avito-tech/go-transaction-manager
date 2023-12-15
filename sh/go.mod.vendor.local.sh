#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

go mod vendor &

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"

    cd $driver && go mod vendor &

    cd $ROOT
  fi
done

wait
