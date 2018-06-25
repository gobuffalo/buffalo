package ci

import (
	"github.com/gobuffalo/buffalo/meta"
)

// Generator for setting CI config in a Buffalo app
type Generator struct {
	App      meta.App
	Provider string
	DBType   string
}

// New CI config generator
func New() Generator {
	return Generator{
		App:      meta.New("."),
		Provider: "travis",
		DBType:   "none",
	}
}
