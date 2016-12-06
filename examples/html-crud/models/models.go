package models

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
