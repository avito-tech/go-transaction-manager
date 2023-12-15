cd ../

ROOT=$(pwd)

go mod tidy

for driver in ./drivers/*/; do
  if [ -d "$driver" ]; then
    echo "\n$driver"
    cd $driver && go mod tidy

    cd $ROOT
  fi
done