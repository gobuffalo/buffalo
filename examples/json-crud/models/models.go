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
	DB, err = pop.Connect(defaults.String(os.Getenv("GO_ENV"), "development"))
	if err != nil {
		log.Fatal(err)
	}
}
