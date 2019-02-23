FROM gobuffalo/buffalo:latest

ARG CODECOV_TOKEN

ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo
RUN rm -rf $BP
RUN mkdir -p $BP
WORKDIR $BP

COPY . .
RUN bash ./it.sh
