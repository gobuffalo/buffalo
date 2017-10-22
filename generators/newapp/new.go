package newapp

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets/standard"
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/buffalo/generators/docker"
	"github.com/gobuffalo/buffalo/generators/refresh"
	"github.com/gobuffalo/buffalo/generators/soda"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

// Run returns a generator to create a new application
func (a Generator) Run(root string, data makr.Data) error {
	g := makr.New()

	if a.Force {
		os.RemoveAll(a.Root)
	}

	g.Add(makr.NewCommand(makr.GoGet("golang.org/x/tools/cmd/goimports", "-u")))
	g.Add(makr.NewCommand(makr.GoGet("github.com/golang/dep/cmd/dep", "-u")))
	g.Add(makr.NewCommand(makr.GoGet("github.com/motemen/gore", "-u")))

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "newapp"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}
	data["name"] = a.Name
	if err := refresh.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	if a.CIProvider == "travis" {
		g.Add(makr.NewFile(".travis.yml", nTravis))
	} else if a.CIProvider == "gitlab-ci" {
		if a.WithPop {
			if a.DBType == "postgres" {
				data["testDbUrl"] = "postgres://postgres:postgres@postgres:5432/" + a.Name.File() + "_test?sslmode=disable"
			} else if a.DBType == "mysql" {
				data["testDbUrl"] = "mysql://root:root@mysql:3306/" + a.Name.File() + "_test"
			} else {
				data["testDbUrl"] = ""
			}
			g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCi))
		} else {
			g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCiNoPop))
		}
	}

	if !a.AsAPI {
		if a.WithWebpack {
			w := webpack.New()
			w.App = a.App
			if err := w.Run(root, data); err != nil {
				return errors.WithStack(err)
			}
		} else {
			if err := standard.Run(root, data); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	if a.WithPop {
		sg := soda.New()
		sg.App = a.App
		sg.Dialect = a.DBType
		data["appPath"] = a.Root
		data["name"] = a.Name.File()
		if err := sg.Run(root, data); err != nil {
			return errors.WithStack(err)
		}
	}
	if a.AsAPI {
		g.Add(makr.Func{
			Runner: func(path string, data makr.Data) error {
				return os.RemoveAll(filepath.Join(path, "templates"))
			},
		})
		g.Add(makr.Func{
			Runner: func(path string, data makr.Data) error {
				return os.RemoveAll(filepath.Join(path, "locales"))
			},
		})
	}
	if a.Docker != "none" {
		o := docker.New()
		o.App = a.App
		o.Version = a.Version
		if err := o.Run(root, data); err != nil {
			return errors.WithStack(err)
		}
	}
	g.Add(makr.NewCommand(a.goGet()))

	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			g.Fmt(root)
			return nil
		},
	})

	if _, err := exec.LookPath("git"); err == nil {
		g.Add(makr.NewCommand(exec.Command("git", "init")))
		g.Add(makr.NewCommand(exec.Command("git", "add", ".")))
		g.Add(makr.NewCommand(exec.Command("git", "commit", "-m", "Initial Commit")))
	}
	data["opts"] = a
	return g.Run(root, data)
}

func (a Generator) goGet() *exec.Cmd {
	if a.WithDep {
		if _, err := exec.LookPath("dep"); err == nil {
			return exec.Command("dep", "init")
		}
	}
	appArgs := []string{"get", "-t"}
	if a.Verbose {
		appArgs = append(appArgs, "-v")
	}
	appArgs = append(appArgs, "./...")
	return exec.Command(envy.Get("GO_BIN", "go"), appArgs...)
}

const nTravis = `
language: go

go:
  - 1.8.x

env:
  - GO_ENV=test

{{ if eq .opts.DBType "postgres" -}}
services:
  - postgresql
{{ end -}}

before_script:
{{ if eq .opts.DBType "postgres" -}}
  - psql -c 'create database {{.opts.Name.File}}_test;' -U postgres
{{ end -}}
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go_import_path: {{.opts.PackagePkg}}

install:
  - go get github.com/gobuffalo/buffalo/buffalo
{{ if .opts.WithDep -}}
  - go get github.com/golang/dep/cmd/dep
  - dep ensure
{{ else -}}
  - go get $(go list ./... | grep -v /vendor/)
{{ end -}}

script: buffalo test
`

const nGitlabCi = `before_script:
  - apt-get update && apt-get install -y postgresql-client mysql-client
  - ln -s /builds /go/src/$(echo "{{.opts.PackagePkg}}" | cut -d "/" -f1)
  - cd /go/src/{{.opts.PackagePkg}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
  - go get -t -v ./...
  - export PATH="$PATH:$GOPATH/bin"

stages:
  - test

.test-vars: &test-vars
  variables:
    GO_ENV: "test"
    POSTGRES_DB: "{{.opts.Name.File}}_test"
    MYSQL_DATABASE: "{{.opts.Name.File}}_test"
    MYSQL_ROOT_PASSWORD: "root"
    TEST_DATABASE_URL: "{{.testDbUrl}}"

# Golang version choice helper
.use-golang-image: &use-golang-latest
  image: golang:latest

.use-golang-image: &use-golang-1-8
  image: golang:1.8

test:latest:
  <<: *use-golang-latest
  <<: *test-vars
  stage: test
  services:
    - mysql:latest
    - postgres:latest
  script:
    - buffalo test

test:1.8:
  <<: *use-golang-1-8
  <<: *test-vars
  stage: test
  services:
    - mysql:latest
    - postgres:latest
  script:
    - buffalo test
`

const nGitlabCiNoPop = `before_script:
  - ln -s /builds /go/src/$(echo "{{.opts.PackagePkg}}" | cut -d "/" -f1)
  - cd /go/src/{{.opts.PackagePkg}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
  - go get -t -v ./...
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

test:latest:
  <<: *use-golang-latest
  <<: *test-vars
  stage: test
  script:
    - buffalo test

test:1.8:
  <<: *use-golang-1-8
  <<: *test-vars
  stage: test
  script:
    - buffalo test
`
