FROM gobuffalo/buffalo:latest

ARG CODECOV_TOKEN

ENV GOPROXY         https://proxy.golang.org
ENV BP              /src/buffalo

RUN rm -rf $BP
RUN mkdir -p $BP

WORKDIR $BP
COPY . .

RUN go mod tidy
RUN go test -tags "sqlite integration_test" -cover -race -v ./...
