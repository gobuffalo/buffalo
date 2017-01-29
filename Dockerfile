FROM golang:latest

RUN go version

RUN apt-get update
RUN curl -sL https://deb.nodesource.com/setup_7.x | bash
RUN apt-get install -y build-essential nodejs
RUN go get -u github.com/golang/lint/golint

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo

RUN mkdir -p $BP
WORKDIR $BP
ADD . .

RUN go get -v -t ./...

RUN go test -race ./...

RUN golint -set_exit_status ./...

RUN go install ./buffalo

WORKDIR $GOPATH/src/
RUN buffalo new --db-type=sqlite3 hello_world
WORKDIR ./hello_world
RUN cat database.yml
RUN go vet -x ./...
RUN buffalo db create -a
RUN buffalo db migrate -e test
RUN buffalo test -race
RUN buffalo g goth facebook twitter linkedin github
RUN buffalo test -race
RUN buffalo build
