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
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

// App is the representation of a new Buffalo application
type App struct {
	RootPath    string
	Name        string
	Force       bool
	Verbose     bool
	SkipPop     bool
	SkipWebpack bool
	SkipYarn    bool
	DBType      string
	CIProvider  string
	API         bool
	WithDep     bool
	Docker      string
}

// Generator returns a generator to create a new application
func (a *App) Generator(data makr.Data) (*makr.Generator, error) {
	g := makr.New()
	g.Add(makr.NewCommand(makr.GoGet("golang.org/x/tools/cmd/goimports", "-u")))
	g.Add(makr.NewCommand(makr.GoGet("github.com/golang/dep/cmd/dep", "-u")))
	g.Add(makr.NewCommand(makr.GoGet("github.com/motemen/gore", "-u")))

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "newapp"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, f := range files {
		g.Add(makr.NewFile(f.WritePath, f.Body))
	}
	rr, err := refresh.New()
	if err != nil {
		return nil, err
	}
	g.Add(rr)

	if data["ciProvider"] == "travis" {
		g.Add(makr.NewFile(".travis.yml", nTravis))
	} else if data["ciProvider"] == "gitlab-ci" {
		if _, ok := data["withPop"]; ok {
			if data["dbType"] == "postgres" {
				data["testDbUrl"] = "postgres://postgres:postgres@postgres:5432/" + data["name"].(string) + "_test?sslmode=disable"
			} else if data["dbType"] == "mysql" {
				data["testDbUrl"] = "mysql://root:root@mysql:3306/" + data["name"].(string) + "_test"
			} else {
				data["testDbUrl"] = ""
			}
			g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCi))
		} else {
			g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCiNoPop))
		}
	}

	if !a.API {
		if a.SkipWebpack {
			wg, err := standard.New(data)
			if err != nil {
				return g, errors.WithStack(err)
			}
			g.Add(wg)
		} else {
			wg, err := webpack.New(data)
			if err != nil {
				return g, errors.WithStack(err)
			}
			g.Add(wg)
		}
	}
	if !a.SkipPop {
		g.Add(newSodaGenerator())
	}
	if a.API {
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
		dg, err := docker.New()
		if err != nil {
			return g, errors.WithStack(err)
		}
		g.Add(dg)
	}
	g.Add(makr.NewCommand(a.goGet()))

	if _, err := exec.LookPath("git"); err == nil {
		g.Add(makr.NewCommand(exec.Command("git", "init")))
		g.Add(makr.NewCommand(exec.Command("git", "add", ".")))
		g.Add(makr.NewCommand(exec.Command("git", "commit", "-m", "Initial Commit")))
	}

	return g, nil
}

func (a App) goGet() *exec.Cmd {
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

{{ if eq .dbType "postgres" -}}
services:
  - postgresql
{{ end -}}

before_script:
{{ if eq .dbType "postgres" -}}
  - psql -c 'create database {{.name}}_test;' -U postgres
{{ end -}}
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go_import_path: {{.packagePath}}

install:
  - go get github.com/gobuffalo/buffalo/buffalo
{{ if .withDep -}}
  - go get github.com/golang/dep/cmd/dep
  - dep ensure
{{ else -}}
  - go get $(go list ./... | grep -v /vendor/)
{{ end -}}

script: buffalo test
`

const nGitlabCi = `before_script:
  - apt-get update && apt-get install -y postgresql-client mysql-client
  - ln -s /builds /go/src/$(echo "{{.packagePath}}" | cut -d "/" -f1)
  - cd /go/src/{{.packagePath}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
  - go get -t -v ./...
  - export PATH="$PATH:$GOPATH/bin"

stages:
  - test

.test-vars: &test-vars
  variables:
    GO_ENV: "test"
    POSTGRES_DB: "{{.name}}_test"
    MYSQL_DATABASE: "{{.name}}_test"
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
  - ln -s /builds /go/src/$(echo "{{.packagePath}}" | cut -d "/" -f1)
  - cd /go/src/{{.packagePath}}
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
