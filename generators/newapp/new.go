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
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
)

// Run returns a generator to create a new application
func (a Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	data["version"] = runtime.Version

	if a.AsAPI {
		defer os.RemoveAll(filepath.Join(a.Root, "templates"))
		defer os.RemoveAll(filepath.Join(a.Root, "locales"))
		defer os.RemoveAll(filepath.Join(a.Root, "public"))
	}
	if a.Force {
		os.RemoveAll(a.Root)
	}

	files, err := generators.FindByBox(Templates)
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

	if sg := a.setupPop(root, data); sg != nil {
		g.Add(sg)
	}

	if err := a.setupDocker(root, data); err != nil {
		return errors.WithStack(err)
	}

	if _, err := exec.LookPath("goimports"); err != nil {
		g.Add(makr.NewCommand(makr.GoGet("golang.org/x/tools/cmd/goimports")))
	}

	if a.WithDep {
		data["addPrune"] = true
		g.Add(makr.NewFile("Gopkg.toml", GopkgTomlTmpl))
		if _, err := exec.LookPath("dep"); err != nil {
			// This step needs to be in a separate generator, because goGet() exec.Command
			// checks if the executable exists (so before running the generator).
			gg := makr.New()
			gg.Add(makr.NewCommand(makr.GoGet("github.com/golang/dep/cmd/dep")))
			if err := gg.Run(root, data); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	for _, c := range a.goGet() {
		g.Add(makr.NewCommand(c))
	}
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
	data["version"] = runtime.Version
	if err := o.Run(root, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a Generator) setupPop(root string, data makr.Data) *makr.Generator {
	if !a.WithPop {
		return nil
	}

	sg := soda.New()
	sg.App = a.App
	sg.Dialect = a.DBType

	g := makr.New()
	for k, v := range data {
		g.Data[k] = v
	}
	g.Data["appPath"] = a.Root
	g.Data["name"] = a.Name.File()
	g.Add(sg)
	return g
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
				data["testDbUrl"] = "mysql://root:root@(mysql:3306)/" + a.Name.File() + "_test?parseTime=true&multiStatements=true&readTimeout=1s"
			} else {
				data["testDbUrl"] = ""
			}
			g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCi))
			break
		}

		g.Add(makr.NewFile(".gitlab-ci.yml", nGitlabCiNoPop))
	}
}

func (a Generator) goGet() []*exec.Cmd {
	cd, _ := os.Getwd()
	defer os.Chdir(cd)
	os.Chdir(a.Root)

	if a.WithDep {
		return []*exec.Cmd{exec.Command("dep", "ensure", "-v")}
	}

	if a.WithModules {
		return a.goGetMod()
	}

	appArgs := []string{"get", "-t"}
	if a.Verbose {
		appArgs = append(appArgs, "-v")
	}
	appArgs = append(appArgs, "./...")
	return []*exec.Cmd{exec.Command(envy.Get("GO_BIN", "go"), appArgs...)}
}

func (a Generator) goGetMod() []*exec.Cmd {
	var cmds []*exec.Cmd
	cmd := exec.Command(envy.Get("GO_BIN", "go"), "get", "github.com/gobuffalo/buffalo@"+runtime.Version)
	if a.Verbose {
		cmd.Args = append(cmd.Args, "-v")
	}
	cmds = append(cmds, cmd)
	cmds = append(cmds, exec.Command(envy.Get("GO_BIN", "go"), "get", "-u", "github.com/gobuffalo/events"))
	cmds = append(cmds, exec.Command(envy.Get("GO_BIN", "go"), "mod", "tidy"))
	return cmds
}

const nTravis = `language: go

go:
	- "1.11.x"

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
{{- if eq .opts.DBType "postgres" }}
  - apt-get update && apt-get install -y postgresql-client
{{- else if eq .opts.DBType "mysql" }}
  - apt-get update && apt-get install -y mysql-client
{{- end }}
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

test:
  # Change to "<<: *use-golang-latest" to use the latest Go version
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

// GopkgTomlTmpl is the default dep Gopkg.toml
const GopkgTomlTmpl = `
[[constraint]]
	name = "github.com/gobuffalo/buffalo"
	{{- if eq .version "development" }}
	branch = "development"
	{{- else }}
	version = "{{.version}}"
	{{- end}}

{{ if .addPrune }}
[prune]
	go-tests = true
	unused-packages = true
{{ end }}

	# DO NOT DELETE
	[[prune.project]] # buffalo
		name = "github.com/gobuffalo/buffalo"
		unused-packages = false

	# DO NOT DELETE
	[[prune.project]] # pop
		name = "github.com/gobuffalo/pop"
		unused-packages = false
`
