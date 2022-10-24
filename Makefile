CVPKG=go list ./... | grep -v mocks | grep -v internal/
GO_TEST=go test `$(CVPKG)` -race
COVERAGE_FILE="coverage.out"

test:
	$(GO_TEST)

test.coverage:
	$(GO_TEST) -covermode=atomic -coverprofile=$(COVERAGE_FILE)

fmt:
	go fmt ./...

lint:
	golangci-lint run -v

generate:
	go generate ./...
