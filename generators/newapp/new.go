package newapp

import (
	"os/exec"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets/standard"
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/buffalo/generators/refresh"
	"github.com/gobuffalo/makr"
)

// App is the representation of a new Buffalo application
type App struct {
	RootPath    string
	Name        string
	Force       bool
	Verbose     bool
	SkipPop     bool
	SkipWebpack bool
	WithYarn    bool
	DBType      string
	CIProvider  string
}

// Generator returns a generator to create a new application
func (a *App) Generator(data makr.Data) (*makr.Generator, error) {
	g := makr.New()
	files, err := generators.Find("newapp")
	if err != nil {
		return nil, err
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
				data["testDbUrl"] = "postgres://postgres:postgres@postgres:5432/" + data["name"].(string) + "_test"
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

	g.Add(makr.NewCommand(makr.GoGet("github.com/motemen/gore")))
	g.Add(makr.NewCommand(makr.GoInstall("github.com/motemen/gore")))
	if a.SkipWebpack {
		wg, err := standard.New(data)
		if err != nil {
			return g, err
		}
		g.Add(wg)
	} else {
		wg, err := webpack.New(data)
		if err != nil {
			return g, err
		}
		g.Add(wg)
	}
	if !a.SkipPop {
		g.Add(newSodaGenerator())
	}
	g.Add(makr.NewCommand(a.goGet()))
	g.Add(makr.NewCommand(makr.GoFmt()))

	return g, nil
}

func (a App) goGet() *exec.Cmd {
	appArgs := []string{"get", "-t"}
	if a.Verbose {
		appArgs = append(appArgs, "-v")
	}
	appArgs = append(appArgs, "./...")
	return exec.Command("go", appArgs...)
}

const nTravis = `language: go
env:
- GO_ENV=test

before_script:
  - psql -c 'create database {{.name}}_test;' -U postgres
	- mysql -e 'CREATE DATABASE {{.name}}_test;'
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go:
  - 1.7.x
  - master

go_import_path: {{ .packagePath }}
`

const nGitlabCi = `before_script:
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
.use-golang-latest: &use-golang-latest
  image: golang:latest

.use-golang-latest: &use-golang-1-7
  image: golang:1.7

test:latest:
  <<: *use-golang-latest
  <<: *test-vars
  stage: test
  services:
    - mysql:latest
    - postgres:latest
  script:
    - buffalo test

test:1.7:
  <<: *use-golang-1-7
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
.use-golang-latest: &use-golang-latest
  image: golang:latest

.use-golang-latest: &use-golang-1-7
  image: golang:1.7

test:latest:
  <<: *use-golang-latest
  <<: *test-vars
  stage: test
  script:
    - buffalo test

test:1.7:
  <<: *use-golang-1-7
  <<: *test-vars
  stage: test
  script:
    - buffalo test
`
