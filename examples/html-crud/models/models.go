package models

import (
	"log"
	"os"
	"time"

	"github.com/markbates/buffalo"
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

var TransactionMW = func(h buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// wrap all requests in a transaction and set the length
		// of time doing things in the db to the log.
		return DB.Transaction(func(tx *pop.Connection) error {
			start := tx.Elapsed
			defer func() {
				finished := tx.Elapsed
				elapsed := time.Duration(finished - start)
				c.LogField("db", elapsed)
			}()
			c.Set("tx", tx)
			return h(c)
		})
	}
}
