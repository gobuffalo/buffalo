FROM gobuffalo/buffalo:latest

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo
RUN rm -rf $BP
RUN mkdir -p $BP
WORKDIR $BP

COPY . .
RUN go mod tidy -v
RUN go test -tags "sqlite integration_test" -cover -race -v ./...
