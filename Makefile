TAGS ?= "sqlite vbuffalo"
INSTALL ?= install -v -tags ${TAGS} ./...

deps:
	go install -v github.com/gobuffalo/packr/packr

install: deps
	packr
	go $(INSTALL)
	packr clean

test: deps
	go test -tags ${TAGS} ./...

vgo-install: deps
	packr
	vgo $(INSTALL)
	packr clean
