package cmd

import "github.com/markbates/gentronics"

func newSodaGenerator() *gentronics.Generator {
	g := gentronics.New()

	f := gentronics.NewFile("models/models.go", nModels)
	f.Should = func(data gentronics.Data) bool {
		if p, ok := data["withPop"]; ok {
			return p.(bool)
		}
		return false
	}
	g.Add(f)
	return g
}

const nModels = `package models

import (
	"log"
	"os"

	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
)

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
