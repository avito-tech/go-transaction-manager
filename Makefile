DIR=$(PWD)

GO_TEST=cd ./sh && bash ./go.test.sh
GO_TEST_COVERAGE=cd ./sh && bash ./go.test.coverage.sh

GO_TEST_WITH_REAL_DB=--tags=with_real_db

test:
	$(GO_TEST)

test.with_real_db:
	$(GO_TEST) $(GO_TEST_WITH_REAL_DB)

test.coverage:
	$(GO_TEST_COVERAGE)

test.coverage.with_real_db:
	$(GO_TEST_COVERAGE) $(GO_TEST_WITH_REAL_DB)

fmt:
	cd sh && sh ./go.fmt.sh

lint:
	cd sh && sh ./lint.sh

lint.verbose:
	cd sh && sh ./lint.sh -v

lint.cache.clean:
	golangci-lint cache clean

generate:
	go generate ./...

go.mod.tidy:
	cd sh && sh ./go.mod.tidy.sh

go.mod.vendor:
	cd sh && sh ./go.mod.vendor.sh

go.work.sync:
	cd sh && sh ./go.work.sync.sh


git.tag: git.tag.create git.tag.push

# 1.0, "v2." added automatically
# make git.tag version="0.0-rc1"
git.tag.create:
	cd sh && sh ./git.tag.sh $(version)

git.tag.push:
	git push origin --tags