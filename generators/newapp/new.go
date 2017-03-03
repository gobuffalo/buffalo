package newapp

import (
	"os/exec"

	"github.com/gobuffalo/buffalo/buffalo/cmd/generate"
	"github.com/gobuffalo/buffalo/generators/common"
	"github.com/gobuffalo/buffalo/generators/refresh"
	"github.com/markbates/gentronics"
	sg "github.com/markbates/pop/soda/cmd/generate"
)

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

func (a *App) Generator(data gentronics.Data) (*gentronics.Generator, error) {
	g := gentronics.New()
	files, err := common.Find("newapp")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		g.Add(gentronics.NewFile(f.WritePath, f.Body))
	}

	g.Add(refresh.New())
	g.Add(gentronics.NewFile(".codeclimate.yml", nCodeClimate))

	if data["ciProvider"] == "travis" {
		g.Add(gentronics.NewFile(".travis.yml", nTravis))
	}

	g.Add(gentronics.NewFile(".gitignore", nGitignore))
	g.Add(gentronics.NewCommand(generate.GoGet("github.com/markbates/refresh/...")))
	g.Add(gentronics.NewCommand(generate.GoInstall("github.com/markbates/refresh")))
	g.Add(gentronics.NewCommand(generate.GoGet("github.com/markbates/grift/...")))
	g.Add(gentronics.NewCommand(generate.GoInstall("github.com/markbates/grift")))
	g.Add(gentronics.NewCommand(generate.GoGet("github.com/motemen/gore")))
	g.Add(gentronics.NewCommand(generate.GoInstall("github.com/motemen/gore")))
	g.Add(generate.NewWebpackGenerator(data))
	g.Add(newSodaGenerator())
	g.Add(gentronics.NewCommand(a.GoGet()))
	g.Add(generate.Fmt)

	return g, nil
}

func (a App) GoGet() *exec.Cmd {
	appArgs := []string{"get", "-t"}
	if a.Verbose {
		appArgs = append(appArgs, "-v")
	}
	appArgs = append(appArgs, "./...")
	return exec.Command("go", appArgs...)
}

const nGitignore = `vendor/
**/*.log
**/*.sqlite
.idea/
bin/
tmp/
node_modules/
.sass-cache/
rice-box.go
public/assets/
{{ name }}
`

const nCodeClimate = `engines:
  fixme:
    enabled: true
  gofmt:
    enabled: true
  golint:
    enabled: true
  govet:
    enabled: true
exclude_paths:
  - grifts/**/*
  - "**/*_test.go"
  - "*_test.go"
  - "**_test.go"
  - logs/*
  - public/*
  - templates/*
ratings:
  paths:
    - "**.go"

`

const nTravis = `language: go
env:
- GO_ENV=test

before_script:
  - psql -c 'create database {{name}}_test;' -U postgres
	- mysql -e 'CREATE DATABASE {{name}}_test;'
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go:
  - 1.7.x
  - master

go_import_path: {{ packagePath }}
`

func newSodaGenerator() *gentronics.Generator {
	g := gentronics.New()

	should := func(data gentronics.Data) bool {
		if _, ok := data["withPop"]; ok {
			return ok
		}
		return false
	}

	f := gentronics.NewFile("models/models.go", nModels)
	f.Should = should
	g.Add(f)

	c := gentronics.NewCommand(generate.GoGet("github.com/markbates/pop/..."))
	c.Should = should
	g.Add(c)

	c = gentronics.NewCommand(generate.GoInstall("github.com/markbates/pop/soda"))
	c.Should = should
	g.Add(c)

	g.Add(&gentronics.Func{
		Should: should,
		Runner: func(rootPath string, data gentronics.Data) error {
			data["dialect"] = data["dbType"]
			return sg.GenerateConfig("./database.yml", data)
		},
	})

	return g
}

const nModels = `package models

import (
	"log"
	"os"

	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	env := defaults.String(os.Getenv("GO_ENV"), "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}
`
