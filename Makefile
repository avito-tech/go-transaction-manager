DIR=$(PWD)

GO_TEST=cd ./sh && bash ./go.test.sh
GO_TEST_COVERAGE=cd ./sh && bash ./go.test.coverage.sh

GO_TEST_WITH_REAL_DB=--tags=with_real_db
# -count=1 disables the test result cache so tests always run
GO_TEST_NO_CACHE=-count=1

test:
	$(GO_TEST) $(GO_TEST_NO_CACHE)

test.with_real_db:
	$(GO_TEST) $(GO_TEST_WITH_REAL_DB) $(GO_TEST_NO_CACHE)

test.coverage:
	$(GO_TEST_COVERAGE) $(GO_TEST_NO_CACHE)

test.coverage.with_real_db:
	$(GO_TEST_COVERAGE) $(GO_TEST_WITH_REAL_DB) $(GO_TEST_NO_CACHE)

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
	go work sync


tag: git.tag tag.pkg

tag.pkg:
	cd sh && sh ./tag.pkg.sh $(version)

git.tag: git.tag.create git.tag.push

# 1.0, "v2." added automatically
# make git.tag version="0.0-rc1"
git.tag.create:
	cd sh && sh ./git.tag.sh $(version)

git.tag.push:
	git push origin --tags
