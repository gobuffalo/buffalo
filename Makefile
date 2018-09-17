TAGS ?= "sqlite"
GO_BIN ?= go

install:
	packr
	$(GO_BIN) install -v .

deps:
	$(GO_BIN) get github.com/gobuffalo/packr/packr
	$(GO_BIN) get -tags ${TAGS} -t ./...

build:
	packr
	$(GO_BIN) build -v .

test:
	packr
	$(GO_BIN) test -tags ${TAGS} ./...

ci-test: deps
	docker build . --no-cache

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u
	$(GO_BIN) mod tidy
	packr
	make test
	make install

release-test:
	make ci-test

release:
	$(GO_BIN) get github.com/gobuffalo/release
	release -y -f runtime/version.go
