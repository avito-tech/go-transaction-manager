cd ../ && go mod tidy

for driver in */../drivers/ ; do
  if [ -d "$driver" ]; then
    cd $driver && go mod vendor
  fi
done