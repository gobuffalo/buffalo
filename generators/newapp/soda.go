package newapp

import (
	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	sg "github.com/markbates/pop/soda/cmd/generate"
)

func newSodaGenerator() *makr.Generator {
	g := makr.New()

	should := func(data makr.Data) bool {
		if _, ok := data["withPop"]; ok {
			return ok
		}
		return false
	}

	f := makr.NewFile("models/models.go", nModels)
	f.Should = should
	g.Add(f)

	c := makr.NewCommand(generators.GoGet("github.com/markbates/pop/..."))
	c.Should = should
	g.Add(c)

	c = makr.NewCommand(generators.GoInstall("github.com/markbates/pop/soda"))
	c.Should = should
	g.Add(c)

	g.Add(&makr.Func{
		Should: should,
		Runner: func(rootPath string, data makr.Data) error {
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
