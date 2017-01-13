FROM golang:latest

RUN go version

RUN apt-get update
RUN curl -sL https://deb.nodesource.com/setup_7.x | bash
RUN apt-get install -y build-essential nodejs

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo

RUN mkdir -p $BP
WORKDIR $BP
ADD . .

RUN go get -v -t github.com/Masterminds/glide
RUN rm -rf vendor/
RUN glide i

RUN go test -race $(glide novendor)

RUN go install ./buffalo

WORKDIR $GOPATH/src/
RUN buffalo new --db-type=sqlite3 hello_world
WORKDIR ./hello_world
RUN cat database.yml
RUN go vet -x $(glide novendor)
RUN buffalo db create -a
RUN buffalo db migrate -e test
RUN buffalo test -race
RUN buffalo g goth facebook twitter linkedin github
RUN buffalo test -race
RUN buffalo build
