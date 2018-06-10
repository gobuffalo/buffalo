TAGS="sqlite vbuffalo"
INSTALL=install -v -tags ${TAGS} ./...

install:
	packr
	go $(INSTALL)
	packr clean

test:
	go test -tags ${TAGS} ./...

vgo-install:
	packr
	vgo $(INSTALL)
	packr clean
