CVPKG=go list ./... | grep -v mocks | grep -v /internal/
GO_TEST=go test `$(CVPKG)` -race
COVERAGE_FILE="coverage.out"
COVERAGE_TMP_FILE="coverage.out.tmp"

test:
	$(GO_TEST)

test.coverage:
	$(GO_TEST) -covermode=atomic -coverprofile=$(COVERAGE_TMP_FILE)

fmt:
	go fmt ./...

lint:
	golangci-lint run -v
