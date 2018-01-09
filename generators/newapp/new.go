package newapp

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	if a.AsAPI {
		defer os.RemoveAll(filepath.Join(a.Root, "templates"))
		defer os.RemoveAll(filepath.Join(a.Root, "locales"))
	}
	if a.Force {
		os.RemoveAll(a.Root)
	}

	g.Add(makr.NewCommand(makr.GoGet("golang.org/x/tools/cmd/goimports", "-u")))
	if a.WithDep {
		g.Add(makr.NewCommand(makr.GoGet("github.com/golang/dep/cmd/dep", "-u")))
	}
	g.Add(makr.NewCommand(makr.GoGet("github.com/motemen/gore", "-u")))

	files, err := generators.Find(filepath.Join(generators.TemplatesPath, "newapp"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		if a.AsAPI {
			if strings.Contains(f.WritePath, "locales") || strings.Contains(f.WritePath, "templates") {
				continue
			}
			g.Add(makr.NewFile(f.WritePath, f.Body))
		} else {
			g.Add(makr.NewFile(f.WritePath, f.Body))
		}

	}
	data["name"] = a.Name
	if err := refresh.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	// Add CI configuration, if requested
	if a.CIProvider == "travis" {
		g.Add(makr.NewFile(".travis.yml", nTravis))
	} else if a.CIProvider == "gitlab-ci" {
		if a.WithPop {
			if a.DBType == "postgres" {
				data["testDbUrl"] = "postgres://postgres:postgres@postgres:5432/" + a.Name.File() + "_test?sslmode=disable"
			} else if a.DBType == "mysql" {
				data["testDbUrl"] = "mysql://root:root@(mysql:3306)/" + a.Name.File() + "_test"
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
			w.Bootstrap = a.Bootstrap
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

	if a.VCS == "git" || a.VCS == "bzr" {
		// Execute git or bzr case (same CLI API)
		if _, err := exec.LookPath(a.VCS); err == nil {
			g.Add(makr.NewCommand(exec.Command(a.VCS, "init")))
			g.Add(makr.NewCommand(exec.Command(a.VCS, "add", ".")))
			g.Add(makr.NewCommand(exec.Command(a.VCS, "commit", "-m", "Initial Commit")))
		}
	}

	data["opts"] = a
	return g.Run(root, data)
}

func (a Generator) goGet() *exec.Cmd {
	cd, _ := os.Getwd()
	defer os.Chdir(cd)
	os.Chdir(a.Root)
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

const nTravis = `language: go

go:
  - 1.8.x

env:
  - GO_ENV=test

{{ if eq .opts.DBType "postgres" -}}
services:
  - postgresql
{{- end }}

before_script:
{{- if eq .opts.DBType "postgres" }}
  - psql -c 'create database {{.opts.Name.File}}_test;' -U postgres
{{- end }}
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go_import_path: {{.opts.PackagePkg}}

install:
  - go get github.com/gobuffalo/buffalo/buffalo
{{- if .opts.WithDep }}
  - go get github.com/golang/dep/cmd/dep
  - dep ensure
{{- else }}
  - go get $(go list ./... | grep -v /vendor/)
{{- end }}

script: buffalo test
`

const nGitlabCi = `before_script:
  - ln -s /builds /go/src/$(echo "{{.opts.PackagePkg}}" | cut -d "/" -f1)
  - cd /go/src/{{.opts.PackagePkg}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
{{- if .opts.WithDep }}
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
    POSTGRES_DB: "{{.opts.Name.File}}_test"
{{- else if eq .opts.DBType "mysql" }}
    MYSQL_DATABASE: "{{.opts.Name.File}}_test"
    MYSQL_ROOT_PASSWORD: "root"
{{- end }}
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
{{- if eq .opts.DBType "mysql" }}
    - mysql:latest
{{- else if eq .opts.DBType "postgres" }}
    - postgres:latest
{{- end }}
  script:
    - buffalo test

test:1.8:
  <<: *use-golang-1-8
  <<: *test-vars
  stage: test
  services:
{{- if eq .opts.DBType "mysql" }}
    - mysql:latest
{{- else if eq .opts.DBType "postgres" }}
    - postgres:latest
{{- end }}
  script:
    - buffalo test
`

const nGitlabCiNoPop = `before_script:
  - ln -s /builds /go/src/$(echo "{{.opts.PackagePkg}}" | cut -d "/" -f1)
  - cd /go/src/{{.opts.PackagePkg}}
  - mkdir -p public/assets
  - go get -u github.com/gobuffalo/buffalo/buffalo
{{- if .opts.WithDep }}
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
