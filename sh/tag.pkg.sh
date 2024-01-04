#!/bin/bash

DIR=$(pwd)
drivers=$($DIR/utils/drivers.sh)

cd ../

tagVersion="v2.$1"

echo "\ntrm"
curl https://proxy.golang.org/github.com/avito-tech/go-transaction-manager/trm/v2/@v/$tagVersion.info

for driver in $drivers; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    curl https://proxy.golang.org/github.com/avito-tech/go-transaction-manager/$driver/v2/@v/$tagVersion.info
  fi
done
