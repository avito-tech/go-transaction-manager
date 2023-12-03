cd ../

ROOT=$(pwd)

golist() {
  go list ./... | grep -v mock | grep -v internal/
}

gotest() {
  cd $driver

  go test $(golist) -race "$@"

  cd $ROOT
}

go test $(golist) -race "$@"

for driver in ./drivers/*; do
  if [ -d "$driver" ]; then
    gotest "$@" &
  fi
done

wait
