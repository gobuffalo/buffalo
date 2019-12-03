TAGS ?= "sqlite"
GO_BIN ?= go

install: tidy
	cd ./buffalo && $(GO_BIN) install -tags ${TAGS} -v
	make tidy

tidy:
	$(GO_BIN) mod tidy

build: tidy
	$(GO_BIN) build -v .

test: tidy
	$(GO_BIN) test -short -tags ${TAGS} -cover ./...
	make tidy

lint:
	golangci-lint --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u -tags ${TAGS}
	make test
	make install
	make tidy

