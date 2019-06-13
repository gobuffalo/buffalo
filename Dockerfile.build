FROM golang:latest
EXPOSE 3000
ENV BP=$GOPATH/src/github.com/gobuffalo/buffalo
RUN go version

RUN curl -sL https://deb.nodesource.com/setup_8.x | bash \
&& sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list' \
&& wget -q https://www.postgresql.org/media/keys/ACCC4CF8.asc -O - | apt-key add - \
&& apt-get update \
&& apt-get install -y -q build-essential nodejs sqlite3 libsqlite3-dev postgresql postgresql-contrib libpq-dev mysql-client vim \
&& rm -rf /var/lib/apt/lists/*

RUN service postgresql start && \
    su -c "psql -c \"ALTER USER postgres  WITH PASSWORD 'postgres';\"" - postgres

RUN go get -u github.com/golang/dep/cmd/dep \
&& go get -tags sqlite -v -u github.com/gobuffalo/pop \
&& go get -tags sqlite -v -u github.com/gobuffalo/buffalo-pop \
&& go get -v -u github.com/gobuffalo/packr/packr \
&& go get -v -u github.com/gobuffalo/packr/v2/packr2 \
&& go get -v -u github.com/markbates/filetest \
&& go get -v -u github.com/markbates/grift \
&& go get -v -u github.com/markbates/refresh \
&& rm -rfv $GOPATH/src && mkdir -p $BP

# Install golangci
RUN  curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.16.0

WORKDIR $BP

ADD go.mod .
ADD go.sum .

# preps the module cache for those using modules in their images
ENV GO111MODULE=on
RUN go mod download
ENV GO111MODULE=off


RUN npm install -g --no-progress yarn \
&& yarn config set yarn-offline-mirror /npm-packages-offline-cache \
&& yarn config set yarn-offline-mirror-pruning true

COPY . .

RUN go get -tags sqlite -t -v ./... && packr2 && go install -v -tags sqlite ./buffalo

RUN buffalo version

WORKDIR $GOPATH/src
