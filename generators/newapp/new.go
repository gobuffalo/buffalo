package newapp

import (
	"fmt"
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
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// Run returns a generator to create a new application
func (a Generator) Run(root string, data makr.Data) error {
	g := makr.New()

	if a.AsAPI {
		defer os.RemoveAll(filepath.Join(a.Root, "templates"))
		defer os.RemoveAll(filepath.Join(a.Root, "locales"))
		defer os.RemoveAll(filepath.Join(a.Root, "public"))
	}
	if a.Force {
		os.RemoveAll(a.Root)
	}

	g.Add(makr.NewCommand(makr.GoGet("golang.org/x/tools/cmd/goimports", "-u")))
	if a.WithDep {
		g.Add(makr.NewCommand(makr.GoGet("github.com/golang/dep/cmd/dep", "-u")))
	}

	files, err := generators.FindByBox(packr.NewBox("../newapp/templates"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range files {
		if !a.AsAPI {
			g.Add(makr.NewFile(f.WritePath, f.Body))
			continue
		}

		if strings.Contains(f.WritePath, "locales") || strings.Contains(f.WritePath, "templates") || strings.Contains(f.WritePath, "public") {
			continue
		}

		g.Add(makr.NewFile(f.WritePath, f.Body))
	}

	data["name"] = a.Name
	if err := refresh.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	a.setupCI(g, data)

	if err := a.setupWebpack(root, data); err != nil {
		return errors.WithStack(err)
	}

	if err := a.setupPop(root, data); err != nil {
		return errors.WithStack(err)
	}

	if err := a.setupDocker(root, data); err != nil {
		return errors.WithStack(err)
	}

	g.Add(makr.NewCommand(a.goGet()))

	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			g.Fmt(root)
			return nil
		},
	})

	a.setupVCS(g)

	data["opts"] = a
	return g.Run(root, data)
}

func (a Generator) setupVCS(g *makr.Generator) {
	if a.VCS != "git" && a.VCS != "bzr" {
		return
	}
	// Execute git or bzr case (same CLI API)
	if _, err := exec.LookPath(a.VCS); err != nil {
		return
	}

	// Create .gitignore or .bzrignore
	g.Add(makr.NewFile(fmt.Sprintf(".%signore", a.VCS), nVCSIgnore))
	g.Add(makr.NewCommand(exec.Command(a.VCS, "init")))
	args := []string{"add", "."}
	if a.VCS == "bzr" {
		// Ensure Bazaar is as quiet as Git
		args = append(args, "-q")
	}
	g.Add(makr.NewCommand(exec.Command(a.VCS, args...)))
	g.Add(makr.NewCommand(exec.Command(a.VCS, "commit", "-q", "-m", "Initial Commit")))
}

func (a Generator) setupDocker(root string, data makr.Data) error {
	if a.Docker == "none" {
		return nil
	}

	o := docker.New()
	o.App = a.App
	o.Version = a.Version
	if err := o.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a Generator) setupPop(root string, data makr.Data) error {
	if !a.WithPop {
		return nil
	}

	sg := soda.New()
	sg.App = a.App
	sg.Dialect = a.DBType
	data["appPath"] = a.Root
	data["name"] = a.Name.File()

	if err := sg.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a Generator) setupWebpack(root string, data makr.Data) error {
	if a.AsAPI {
		return nil
	}

	if a.WithWebpack {
		w := webpack.New()
		w.App = a.App
		w.Bootstrap = a.Bootstrap
		if err := w.Run(root, data); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	if err := standard.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a Generator) setupCI(g *makr.Generator, data makr.Data) {

	switch a.CIProvider {
	case "travis":
		g.Add(makr.NewFile(".travis.yml", nTravis))
	case "gitlab-ci":
		if a.WithPop {
			if a.DBType == "postgres" {
				data["testDbUrl"] = "postgres://postgres:postgres@postgres:5432/" + a.Name.File() + "_test?sslmode=disable"
			} else if a.DBType == "mysql" {
				data["testDbUrl"] = "mysql://root:root@(mysql:3306)/" + a.Name.File() + "_test"
			} else {
				data["testDbUrl"] = ""
			}
			g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCi))
			break
		}

		g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCiNoPop))
	}
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

const nVCSIgnore = `vendor/
**/*.log
**/*.sqlite
.idea/
bin/
tmp/
node_modules/
.sass-cache/
*-packr.go
public/assets/
{{ .opts.Name.File }}
.vscode/
.grifter/
.env
`
