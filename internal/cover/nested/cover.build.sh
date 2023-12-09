cd ./module1/cmd

mkdir -p codecov

go build -cover -o run.exe .
GOCOVERDIR=./codecov ./run.exe

# go test ./... -covermode=atomic -coverprofile=coverage.out

cd ../../module2/cmd

mkdir -p codecov

go build -cover -o run.exe .
GOCOVERDIR=./codecov ./run.exe

# go test ./... -covermode=atomic -coverprofile=coverage.out

cd ../../

go tool covdata textfmt -i=module1/cmd/codecov,module2/cmd/codecov -o ./output.out