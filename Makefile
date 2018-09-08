TAGS ?= "sqlite"
INSTALL ?= install -v -tags ${TAGS} ./...

GO_BIN ?= go
GO_GET ?= $(GO_BIN) get -tags "sqlite" -v -t github.com/gobuffalo/buffalo/...

install: deps
	packr
	$(GO_GET)
	$(GO_BIN) $(INSTALL)
	packr clean

ifeq ("$(GO_BIN)","vgo")
	GO_GET = vgo version
endif

deps:
	$(GO_BIN) install -v github.com/gobuffalo/packr/packr

test:
	$(GO_BIN) test -tags ${TAGS} ./...

ci-test:
	$(GO_BIN) test -tags ${TAGS} -race -v ./...
	docker build .

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u
	$(GO_BIN) mod tidy
	packr
	make test
