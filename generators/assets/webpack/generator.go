package webpack

import "github.com/gobuffalo/buffalo/meta"

// Generator for creating a new webpack setup
type Generator struct {
	meta.App
}

// New webpack generator
func New() Generator {
	return Generator{
		App: meta.New("."),
	}
}
