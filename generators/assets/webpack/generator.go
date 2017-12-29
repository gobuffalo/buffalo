package webpack

import "github.com/gobuffalo/buffalo/meta"

// Generator for creating a new webpack setup
type Generator struct {
	meta.App
	Bootstrap int `json:"bootstrap"`
}

// New webpack generator
func New() Generator {
	return Generator{
		App: meta.New("."),
	}
}
