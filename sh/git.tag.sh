#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

ROOT=$(pwd)

tagVersion="v2.$1"

echo "\ntrm"
git tag trm/$tagVersion

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    git tag $driver/$tagVersion
  fi
done
