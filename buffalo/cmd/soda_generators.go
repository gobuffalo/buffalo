package cmd

import (
	"github.com/markbates/gentronics"
	"github.com/markbates/pop/soda/cmd/generate"
)

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

	c := gentronics.NewCommand(goGet("github.com/markbates/pop/..."))
	c.Should = should
	g.Add(c)

	c = gentronics.NewCommand(goInstall("github.com/markbates/pop/soda"))
	c.Should = should
	g.Add(c)

	g.Add(&gentronics.Func{
		Should: should,
		Runner: func(rootPath string, data gentronics.Data) error {
			data["dialect"] = data["dbType"]
			generate.GenerateConfig("./database.yml", data)
			return nil
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
