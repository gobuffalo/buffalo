package newapp

import (
	"os/exec"

	"github.com/gobuffalo/buffalo/buffalo/cmd/generate"
	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/generators/assets/standard"
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/buffalo/generators/refresh"
	"github.com/markbates/gentronics"
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
func (a *App) Generator(data gentronics.Data) (*gentronics.Generator, error) {
	g := gentronics.New()
	files, err := generators.Find("newapp")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		g.Add(gentronics.NewFile(f.WritePath, f.Body))
	}
	rr, err := refresh.New()
	if err != nil {
		return nil, err
	}
	g.Add(rr)

	if data["ciProvider"] == "travis" {
		g.Add(gentronics.NewFile(".travis.yml", nTravis))
	}

	g.Add(gentronics.NewCommand(generators.GoGet("github.com/markbates/refresh/...")))
	g.Add(gentronics.NewCommand(generators.GoInstall("github.com/markbates/refresh")))
	g.Add(gentronics.NewCommand(generators.GoGet("github.com/markbates/grift/...")))
	g.Add(gentronics.NewCommand(generators.GoInstall("github.com/markbates/grift")))
	g.Add(gentronics.NewCommand(generators.GoGet("github.com/motemen/gore")))
	g.Add(gentronics.NewCommand(generators.GoInstall("github.com/motemen/gore")))
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
	g.Add(newSodaGenerator())
	g.Add(gentronics.NewCommand(a.goGet()))
	g.Add(generate.Fmt)

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
  - psql -c 'create database {{name}}_test;' -U postgres
	- mysql -e 'CREATE DATABASE {{name}}_test;'
  - mkdir -p $TRAVIS_BUILD_DIR/public/assets

go:
  - 1.7.x
  - master

go_import_path: {{ packagePath }}
`
