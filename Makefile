TAGS ?= "sqlite"
GO_BIN ?= go

install:
	packr
	$(GO_BIN) install -tags ${TAGS} -v ./buffalo

deps:
	$(GO_BIN) get github.com/gobuffalo/packr/packr
	$(GO_BIN) get -tags ${TAGS} -t ./...

build:
	packr
	$(GO_BIN) build -v .

test:
	packr
	$(GO_BIN) test -tags ${TAGS} ./...

ci-deps:
	$(GO_BIN) get github.com/gobuffalo/packr/packr
	$(GO_BIN) get -tags ${TAGS} -t -u -v ./...

ci-test: ci-deps
	docker build . --no-cache

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u
	packr
	make test
	make install
	$(GO_BIN) mod tidy

release-test:
	make test

release:
	$(GO_BIN) get github.com/gobuffalo/release
	release -y -f runtime/version.go
