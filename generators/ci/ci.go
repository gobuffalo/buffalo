package ci

import (
	"fmt"

	"github.com/gobuffalo/makr"
)

// Run the CI config generator
func (cg Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	should := func(data makr.Data) bool {
		return cg.Provider != "none"
	}
	data["opts"] = cg
	var f makr.File
	switch cg.Provider {
	case "travis":
		f = makr.NewFile(".travis.yml", nTravis)
		f.Should = should
		g.Add(f)
	case "gitlab-ci":
		if cg.DBType != "none" {
			if cg.DBType == "postgres" {
				data["testDbUrl"] = "postgres://postgres:postgres@postgres:5432/" + cg.App.Name.File() + "_test?sslmode=disable"
			} else if cg.DBType == "mysql" {
				data["testDbUrl"] = "mysql://root:root@(mysql:3306)/" + cg.App.Name.File() + "_test?parseTime=true&multiStatements=true&readTimeout=1s"
			} else {
				data["testDbUrl"] = ""
			}
			f = makr.NewFile(".gitlab-ci.yml", nGitlabCi)
			f.Should = should
			g.Add(f)
			break
		}

		f = makr.NewFile(".gitlab-ci.yml", nGitlabCiNoPop)
		f.Should = should
		g.Add(f)
	default:
		return fmt.Errorf("unsupported CI provider %s", cg.Provider)
	}
	return g.Run(root, data)
}

const nTravis = `language: go

go:
  - 1.8.x

env:
  - GO_ENV=test

{{ if eq .opts.DBType "postgres" -}}
services:
  - postgresql
{{- else if eq .opts.DBType "mysql" -}}
services:
  - mysql
{{- end }}

before_script:
{{- if eq .opts.DBType "postgres" }}
  - psql -c 'CREATE DATABASE {{.opts.App.Name.File}}_test;' -U postgres
{{- else if eq .opts.DBType "mysql" }}
  - mysql -e 'CREATE DATABASE {{.opts.App.Name.File}}_test;'
{{- end }}
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go_import_path: {{.opts.App.PackagePkg}}

install:
  - go get github.com/gobuffalo/buffalo/buffalo
{{- if .opts.App.WithDep }}
  - go get github.com/golang/dep/cmd/dep
  - dep ensure
{{- else }}
  - go get $(go list ./... | grep -v /vendor/)
{{- end }}

script: buffalo test
`

const nGitlabCi = `before_script:
{{- if eq .opts.DBType "postgres" }}
	- apt-get update && apt-get install -y postgresql-client
{{- else if eq .opts.DBType "mysql" }}
  - apt-get update && apt-get install -y mysql-client
{{- end }}
  - ln -s /builds /go/src/$(echo "{{.opts.App.PackagePkg}}" | cut -d "/" -f1)
  - cd /go/src/{{.opts.App.PackagePkg}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
{{- if .opts.App.WithDep }}
  - go get github.com/golang/dep/cmd/dep
  - dep ensure
{{- else }}
  - go get -t -v ./...
{{- end }}
  - export PATH="$PATH:$GOPATH/bin"

stages:
  - test

.test-vars: &test-vars
  variables:
    GO_ENV: "test"
{{- if eq .opts.DBType "postgres" }}
    POSTGRES_DB: "{{.opts.App.Name.File}}_test"
{{- else if eq .opts.DBType "mysql" }}
    MYSQL_DATABASE: "{{.opts.App.Name.File}}_test"
    MYSQL_ROOT_PASSWORD: "root"
{{- end }}
    TEST_DATABASE_URL: "{{.testDbUrl}}"

# Golang version choice helper
.use-golang-image: &use-golang-latest
  image: golang:latest

.use-golang-image: &use-golang-1-8
  image: golang:1.8

test:
  # Change to "<<: *use-golang-latest" to use the latest Go version
  <<: *use-golang-1-8
  <<: *test-vars
  stage: test
  services:
{{- if eq .opts.DBType "mysql" }}
    - mysql:5
{{- else if eq .opts.DBType "postgres" }}
    - postgres:latest
{{- end }}
  script:
    - buffalo test
`

const nGitlabCiNoPop = `before_script:
  - ln -s /builds /go/src/$(echo "{{.opts.App.PackagePkg}}" | cut -d "/" -f1)
  - cd /go/src/{{.opts.App.PackagePkg}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
{{- if .opts.App.WithDep }}
  - go get github.com/golang/dep/cmd/dep
  - dep ensure
{{- else }}
  - go get -t -v ./...
{{- end }}
  - export PATH="$PATH:$GOPATH/bin"

stages:
  - test

.test-vars: &test-vars
  variables:
    GO_ENV: "test"

# Golang version choice helper
.use-golang-image: &use-golang-latest
  image: golang:latest

.use-golang-image: &use-golang-1-8
  image: golang:1.8

test:
  # Change to "<<: *use-golang-latest" to use the latest Go version
  <<: *use-golang-1-8
  <<: *test-vars
  stage: test
  script:
    - buffalo test
`
