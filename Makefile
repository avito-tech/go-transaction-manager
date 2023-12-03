GO_TEST=cd ./sh && sh ./go.test.sh
GO_TEST_WITH_REAL_DB=$(GO_TEST) --tags=with_real_db

DIR=$(PWD)
COVERAGE_FILE=`echo $(DIR)/coverage.out`

test:
	$(GO_TEST)

test.with_real_db:
	$(GO_TEST_WITH_REAL_DB)

# TODO see in https://gist.github.com/skarllot/13ebe8220822bc19494c8b076aabe9fc
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


go.mod.tidy:
	cd sh && sh ./go.mod.tidy.sh

go.mod.vendor:
	cd sh && sh ./go.mod.vendor.sh
