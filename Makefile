TAGS ?= "sqlite"
GO_BIN ?= go

install: deps
	packr
	$(GO_BIN) install -tags ${TAGS} -v ./buffalo
	make tidy

tidy:
ifeq ($(GO111MODULE),on)
	$(GO_BIN) mod tidy
else
	echo skipping go mod tidy
endif

deps:
	$(GO_BIN) get github.com/gobuffalo/release
	$(GO_BIN) get github.com/gobuffalo/packr/packr
	$(GO_BIN) get -tags ${TAGS} -t ./...
	make tidy

build:
	packr
	$(GO_BIN) build -v .
	make tidy

test:
	packr
	$(GO_BIN) test -tags ${TAGS} ./...
	make tidy

ci-deps:
	$(GO_BIN) get github.com/gobuffalo/packr/packr
	$(GO_BIN) get -tags ${TAGS} -t -u -v ./...
	make tidy

ci-test:
	docker build . --no-cache

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u -tags ${TAGS}
	make tidy
	packr
	make test
	make install
	make tidy

release-test:
	make test
	make tidy

release:
	make tidy
	release -y -f ./runtime/version.go
	make tidy
