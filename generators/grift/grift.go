package grift

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
)

//New allows to create a new grift task generator
func New(data makr.Data) (*makr.Generator, error) {
	g := makr.New()

	files, err := generators.Find("grift")
	if err != nil {
		return nil, err
	}

	path := filepath.Join("grifts", data["filename"].(string))
	file := files[0]
	g.Add(makr.NewFile(path, file.Body))

	return g, nil
}
