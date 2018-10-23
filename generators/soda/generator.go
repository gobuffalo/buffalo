package soda

import "github.com/gobuffalo/meta"

// Generator for setting soda in a Buffalo app
type Generator struct {
	App     meta.App
	Dialect string
}

// New soda generator
func New() Generator {
	return Generator{
		App:     meta.New("."),
		Dialect: "postgres",
	}
}
