package soda

import (
	"github.com/gobuffalo/makr"
	sg "github.com/gobuffalo/pop/soda/cmd/generate"
)

// Run the soda generator
func (sd Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	defer g.Fmt(root)

	should := func(data makr.Data) bool {
		return sd.App.WithPop
	}

	f := makr.NewFile("models/models.go", nModels)
	f.Should = should
	g.Add(f)

	f = makr.NewFile("models/models_test.go", nModelsTest)
	f.Should = should
	g.Add(f)

	f = makr.NewFile("grifts/db.go", nSeedGrift)
	f.Should = should
	g.Add(f)

	c := makr.NewCommand(makr.GoGet("github.com/gobuffalo/pop/..."))
	c.Should = should
	g.Add(c)

	g.Add(&makr.Func{
		Should: should,
		Runner: func(rootPath string, data makr.Data) error {
			data["dialect"] = sd.Dialect
			return sg.GenerateConfig("./database.yml", data)
		},
	})

	return g.Run(root, data)
}

const nModels = `package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
}
`

const nModelsTest = `package models_test

import (
	"testing"

	"github.com/gobuffalo/suite"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	as := &ModelSuite{suite.NewModel()}
	suite.Run(t, as)
}`

const nSeedGrift = `package grifts

import (
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		return nil
	})

})`
