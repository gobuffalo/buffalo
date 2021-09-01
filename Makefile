TAGS ?= "sqlite"
GO_BIN ?= go

install: deps
	make tidy
	$(GO_BIN) install -tags ${TAGS} -v ./buffalo
	make tidy

tidy:
	$(GO_BIN) mod tidy

deps:
	$(GO_BIN) get github.com/gobuffalo/release
	make tidy

build:
	$(GO_BIN) build -v .
	make tidy

test:
	packr2
	$(GO_BIN) test -tags ${TAGS} -cover ./...
	packr2
	make tidy

ci-deps:
	$(GO_BIN) get github.com/gobuffalo/buffalo-pop
	$(GO_BIN) get -tags ${TAGS} -t -v ./...
	make tidy

ci-test:
	docker build . --no-cache --build-arg TRAVIS_BRANCH=$$(git symbolic-ref --short HEAD)

lint:
	golangci-lint --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u -tags ${TAGS}
	make tidy
	make test
	make install
	make tidy

release-test:
	make tidy

release:
	make tidy
	release -y -f ./runtime/version.go --skip-packr
	make tidy
