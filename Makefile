# GC is the go compiler.
GC = go

export GO111MODULE=on

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	chmod +rwx scripts/*.sh
	scripts/add-license.sh

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	chmod +rwx scripts/*.sh
	scripts/build-storage.sh

# fmt runs the formatter.
.PHONY: fmt
fmt:
	gofumpt -l -w .

# lint runs the linter.
.PHONY: lint
lint:
	golangci-lint run --fix;
	golangci-lint run --config .golangci.yml;

# test runs all of the unit tests locally. Each test is run 5 times to minimize flakiness.
.PHONY: tests
tests:
	$(GC) clean -testcache
	go test -v -count=5 ./...

# ci runs the tests with a ci container
.PHONY: ci
ci:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-ci-tests.sh
