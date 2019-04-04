TAGS ?= "sqlite"
GO_BIN ?= go

install: deps
	make tidy
	$(GO_BIN) install -tags ${TAGS} -v ./buffalo
	make tidy

tidy:
	packr2
ifeq ($(GO111MODULE),on)
	$(GO_BIN) mod tidy
else
	echo skipping go mod tidy
endif

deps:
	$(GO_BIN) get github.com/gobuffalo/release
	$(GO_BIN) get github.com/gobuffalo/packr/v2/packr2
	packr2 clean
ifneq ($(GO111MODULE),on)
	$(GO_BIN) get -tags ${TAGS} -u -t ./...
endif
	make tidy

build:
	packr2
	$(GO_BIN) build -v .
	make tidy

test:
	packr2
	$(GO_BIN) test -tags ${TAGS} ./...
	make tidy

ci-deps:
	$(GO_BIN) get -u github.com/gobuffalo/packr/v2/packr2
	$(GO_BIN) get github.com/gobuffalo/buffalo-pop
	$(GO_BIN) get -tags ${TAGS} -t -v ./...
	make tidy

ci-test:
	docker build . --no-cache --build-arg TRAVIS_BRANCH=$$(git symbolic-ref --short HEAD)

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u -tags ${TAGS}
	make tidy
	packr2
	make test
	make install
	make tidy

release-test:
	make test
	make tidy

release:
	make tidy
	release -y -f ./runtime/version.go --skip-packr
	make tidy
