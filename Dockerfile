FROM golang

ENV BP=$GOPATH/src/github.com/markbates/buffalo

RUN mkdir -p $BP
WORKDIR $BP
ADD . .

RUN go get -v -t ./...

RUN go test -v -race ./...

RUN go install ./buffalo

WORKDIR $GOPATH/src/
RUN buffalo new --db-type=sqlite3 hello_world
WORKDIR ./hello_world
RUN go vet -x ./...
RUN buffalo test -v -race
