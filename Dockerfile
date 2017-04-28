FROM golang:latest

RUN go version

RUN apt-get update
RUN curl -sL https://deb.nodesource.com/setup_7.x | bash
RUN apt-get install -y build-essential nodejs
RUN apt-get install -y sqlite3 libsqlite3-dev

RUN go get -u github.com/golang/lint/golint
RUN go get -u github.com/markbates/filetest

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo

RUN mkdir -p $BP
WORKDIR $BP
ADD . .

RUN go get -v -t ./...

RUN go test -race ./...

RUN golint -set_exit_status ./...

RUN go install ./buffalo

WORKDIR $GOPATH/src/
RUN buffalo new --db-type=sqlite3 hello_world --ci-provider=travis
WORKDIR ./hello_world

RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/new_travis.json

RUN go vet -x ./...
RUN buffalo db create -a
RUN buffalo db migrate -e test
RUN buffalo test -race

RUN buffalo g goth facebook twitter linkedin github
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/goth.json

RUN buffalo g resource admins --skip-model
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_skip_model.json
RUN rm actions/admins_test.go

RUN buffalo test -race
RUN buffalo build

RUN buffalo g resource users name:text email:text
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_model_migration.json

RUN rm models/user_test.go
RUN rm models/user.go
RUN rm actions/users_test.go
RUN rm -rv templates/users

RUN buffalo g resource --type=json users name:text email:text
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_json-xml.json

RUN rm models/user_test.go
RUN rm models/user.go
RUN rm actions/users_test.go

RUN buffalo g resource --type=xml users name:text email:text
RUN filetest -c $GOPATH/src/github.com/gobuffalo/buffalo/buffalo/cmd/filetests/resource_json-xml.json

RUN rm models/user_test.go
RUN rm models/user.go
RUN rm actions/users_test.go

WORKDIR $GOPATH/src
RUN buffalo new --skip-pop simple_world
WORKDIR ./simple_world
RUN buffalo build
