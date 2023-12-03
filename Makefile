CVPKG=go list ./... | grep -v mocks | grep -v internal/
GO_TEST=go test `$(CVPKG)` -race
GO_TEST_WITH_REAL_DB=$(GO_TEST) --tags=with_real_db
COVERAGE_FILE="coverage.out"

test:
	$(GO_TEST)

test.with_real_db:
	$(GO_TEST_WITH_REAL_DB)

test.coverage:
	$(GO_TEST) -covermode=atomic -coverprofile=$(COVERAGE_FILE)

test.coverage.with_real_db:
	$(GO_TEST_WITH_REAL_DB) -covermode=atomic -coverprofile=$(COVERAGE_FILE)

fmt:
	go fmt ./...

lint:
	golangci-lint run -v --timeout=2m

generate:
	go generate ./...


go.update: go.tidy go.vendor

go.prepare: go.vendor

go.tidy:
	./sh/go.mod.tidy.sh

go.vendor:
	./sh/go.mod.vendor.sh
